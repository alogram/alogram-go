// Copyright (c) 2025 Alogram Inc.
// All rights reserved.

package alogram

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alogram/alogram-go/internal/payrisk_v1"
)

func buildValidRequest() payrisk_v1.CheckRequest {
	return payrisk_v1.CheckRequest{
		EventType: payrisk_v1.PtrString("purchase"),
		Entities: payrisk_v1.EntityIds{
			TenantId: payrisk_v1.PtrString("tid_test"),
		},
		Purchase: payrisk_v1.Purchase{
			Amount:   10.0,
			Currency: "USD",
			PaymentMethod: payrisk_v1.PaymentMethod{
				Card: &payrisk_v1.Card{
					Type:        "card",
					Bin:         payrisk_v1.PtrString("424242"),
					CardNetwork: payrisk_v1.CARDNETWORKENUM_VISA.Ptr(),
				},
			},
		},
	}
}

func TestDualTrustInitialization(t *testing.T) {
	// Risk client should block pk_
	_, err := NewAlogramRiskClient(ClientOptions{APIKey: "pk_test"})
	if err == nil || !strings.Contains(err.Error(), "Publishable Key") {
		t.Errorf("Expected RiskClient to block pk_ key, got %v", err)
	}

	// Public client should block sk_
	_, err = NewAlogramPublicClient(ClientOptions{APIKey: "sk_test"})
	if err == nil || !strings.Contains(err.Error(), "Secret Key") {
		t.Errorf("Expected PublicClient to block sk_ key, got %v", err)
	}
}

func TestCheckRiskSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "sk_test" {
			t.Errorf("Expected sk_test key, got %s", r.Header.Get("x-api-key"))
		}

		resp := payrisk_v1.DecisionResponse{
			Decision:  "approve",
			RiskScore: 0.1,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client, _ := NewAlogramRiskClient(ClientOptions{
		BaseURL: ts.URL,
		APIKey:  "sk_test",
	})

	req := buildValidRequest()
	resp, err := client.CheckRisk(context.Background(), req, "", "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.Decision != "approve" {
		t.Errorf("Expected approve, got %s", resp.Decision)
	}
}

func TestPublicClientIngestOnly(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	client, _ := NewAlogramPublicClient(ClientOptions{
		BaseURL: ts.URL,
		APIKey:  "pk_test",
	})

	req := payrisk_v1.SignalsRequest{
		SignalsInteractionVariant: &payrisk_v1.SignalsInteractionVariant{
			SignalType: "interaction",
		},
	}

	err := client.IngestSignals(context.Background(), req, "", "")
	if err != nil {
		t.Errorf("Expected success, got %v", err)
	}
}
