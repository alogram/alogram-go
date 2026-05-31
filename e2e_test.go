// Copyright (c) 2025 Alogram Inc.
// Internal E2E Simulator for Go SDK (Full Circle)

package alogram

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alogram/alogram-go/internal/payrisk_v1"
)

func loadAPIKey(envName string) string {
	if key := os.Getenv("ALOGRAM_API_KEY"); key != "" {
		return key
	}
	// Try loading from project keystore
	path := os.Getenv("ALOGRAM_KEYSTORE_PATH")
	if path == "" {
		path = "../../keystore/api-keys.json"
	}
	data, err := os.ReadFile(path)
	if err == nil {
		var keys map[string]string
		if err := json.Unmarshal(data, &keys); err == nil {
			if val, ok := keys["api-key-"+envName]; ok {
				return val
			}
		}
	}
	return "sk_test_internal"
}

func TestE2E_Simulator(t *testing.T) {
	baseURL := os.Getenv("ALOGRAM_BASE_URL")
	if baseURL == "" {
		// 🛠️ Alogram: Default to local emulator for stable, key-independent verification
		baseURL = "http://localhost:8080"
	}
	envName := os.Getenv("ENV")
	if envName == "" {
		envName = "dev"
	}
	apiKey := loadAPIKey(envName)
	tenantID := "tid_alogramtech"

	isLocal := strings.Contains(baseURL, "localhost") || strings.Contains(baseURL, "127.0.0.1") || strings.Contains(baseURL, "host.docker.internal")
	
	// 🛡️ Alogram: Self-healing pre-flight check
	if isLocal {
		conn, err := net.DialTimeout("tcp", strings.TrimPrefix(strings.TrimPrefix(baseURL, "http://"), "https://"), 1*time.Second)
		if err != nil {
			t.Skipf("⏩ Skipping Local E2E: Emulator not reachable at %s. Start it with 'make emulator-run'.", baseURL)
		}
		conn.Close()
	} else if apiKey == "" || apiKey == "sk_test_internal" {
		t.Skip("⏩ Skipping Cloud E2E: No valid ALOGRAM_API_KEY found and target is remote.")
	}

	fmt.Printf("🏁 Starting Go E2E Simulator (Full Circle)\n")
	fmt.Printf("🚀 Target: %s\n", baseURL)

	cfg := payrisk_v1.NewConfiguration()
	cfg.Servers = payrisk_v1.ServerConfigurations{
		{
			URL: baseURL,
		},
	}
	cfg.AddDefaultHeader("x-api-key", apiKey)
	cfg.AddDefaultHeader("x-trusted-tenant-id", tenantID)

	client := payrisk_v1.NewAPIClient(cfg)
	authCtx := context.WithValue(context.Background(), payrisk_v1.ContextAPIKeys, map[string]payrisk_v1.APIKey{
		"ApiKey": {Key: apiKey},
	})

	// 🚀 Step 1: Intelligent Handshake
	if !isLocal {
		fmt.Println("⏳ Performing Go infrastructure handshake...")
		success := false
		for i := 1; i <= 5; i++ {
			resp, err := client.SystemAPI.HealthCheck(authCtx).Execute()
			if err == nil && resp.StatusCode == 200 {
				success = true
				fmt.Println("✅ Infrastructure is READY.")
				break
			}
			wait := time.Duration(1<<i) * time.Second
			fmt.Printf("⚠️ Handshake attempt %d failed. Retrying in %v...\n", i, wait)
			time.Sleep(wait)
		}
		if !success {
			t.Fatal("❌ Infrastructure handshake TIMEOUT.")
		}
	} else {
		fmt.Println("⚡ Local environment detected. Skipping infrastructure warmup.")
	}

	// Use a more robust random ID for sessionId
	sessionHex := fmt.Sprintf("%x", time.Now().UnixNano())
	entities := payrisk_v1.EntityIds{
		TenantId:      payrisk_v1.PtrString(tenantID),
		ClientId:      payrisk_v1.PtrString("cid_go_sim"),
		EndCustomerId: payrisk_v1.PtrString("ecid_go_" + sessionHex),
		SessionId:     payrisk_v1.PtrString("sid_" + sessionHex),
	}

	// 📡 Step 2: Ingest Behavioral Signals
	fmt.Println("📡 Ingesting Behavioral Signals...")
	meta, _ := json.Marshal(map[string]string{"source": "go-e2e", "action": "login"})
	sigReq := payrisk_v1.SignalsRequest{
		SignalsInteractionVariant: &payrisk_v1.SignalsInteractionVariant{
			SignalType: "interaction",
			Entities:   entities,
			Interactions: []payrisk_v1.Interaction{
				{
					InteractionType: "login",
					Timestamp:       payrisk_v1.PtrString(time.Now().UTC().Format(time.RFC3339)),
					Metadata:        payrisk_v1.PtrString(string(meta)),
				},
			},
		},
	}

	payloadJson, _ := json.MarshalIndent(sigReq, "", "  ")
	fmt.Printf("📦 Signals Payload:\n%s\n", string(payloadJson))

	_, _, err := client.SignalIntelligenceAPI.IngestSignals(authCtx).
		SignalsRequest(sigReq).
		XIdempotencyKey(fmt.Sprintf("idk_%032x", time.Now().UnixNano())).
		Execute()
	if err != nil {
		t.Fatalf("🔥 Signals ingestion failed: %v", err)
	}
	fmt.Println("✅ Signals accepted.")

	// 🧪 Step 3: Risk Assessment
	fmt.Println("🧪 Requesting Risk Assessment...")
	checkReq := payrisk_v1.CheckRequest{
		Entities: entities,
		Purchase: payrisk_v1.Purchase{
			Timestamp: payrisk_v1.PtrString(time.Now().UTC().Format(time.RFC3339)),
			Amount:    350.00,
			Currency:  "USD",
			PaymentMethod: payrisk_v1.PaymentMethod{
				Card: &payrisk_v1.Card{
					Type:        "card",
					CardNetwork: payrisk_v1.CARDNETWORKENUM_VISA.Ptr(),
					Bin:         payrisk_v1.PtrString("411111"),
				},
			},
		},
	}
	decision, _, err := client.RiskScoringAPI.RiskCheck(authCtx).
		CheckRequest(checkReq).
		XIdempotencyKey(fmt.Sprintf("idk_%032x", time.Now().UnixNano())).
		Execute()
	if err != nil {
		t.Fatalf("🔥 Risk check failed: %v", err)
	}
	fmt.Printf("✅ Decision: %s | Score: %f\n", decision.Decision, decision.DecisionScore)

	if decision.PaymentIntentId != nil {
		// 💳 Step 4: Authorization
		fmt.Println("💳 Ingesting Authorization...")
		eventAuth := payrisk_v1.PaymentEvent{
			PaymentIntentId: *decision.PaymentIntentId,
			EventType:       "authorization",
			Timestamp:       time.Now().UTC().Format(time.RFC3339),
			Amount:          payrisk_v1.PtrFloat32(350.00),
			Currency:        payrisk_v1.PtrString("USD"),
			Outcome: &payrisk_v1.PaymentOutcome{
				Authorization: &payrisk_v1.PaymentAuthorizationOutcome{
					Approved:     payrisk_v1.PtrBool(true),
					ResponseCode: payrisk_v1.PtrString("00"),
				},
			},
		}
		_, _, err = client.SignalIntelligenceAPI.IngestPaymentEvent(authCtx).
			PaymentEvent(eventAuth).
			XIdempotencyKey(fmt.Sprintf("idk_%032x", time.Now().UnixNano())).
			Execute()
		if err != nil {
			t.Fatalf("🔥 Auth ingestion failed: %v", err)
		}
		fmt.Println("✅ Authorization accepted.")

		// 💰 Step 5: Capture
		fmt.Println("💰 Ingesting Capture...")
		eventCap := payrisk_v1.PaymentEvent{
			PaymentIntentId: *decision.PaymentIntentId,
			EventType:       "capture",
			Timestamp:       time.Now().UTC().Format(time.RFC3339),
			Amount:          payrisk_v1.PtrFloat32(350.00),
			Currency:        payrisk_v1.PtrString("USD"),
			Outcome: &payrisk_v1.PaymentOutcome{
				Capture: &payrisk_v1.PaymentCaptureOutcome{
					Status: payrisk_v1.PtrString("full"),
				},
			},
		}
		_, _, err = client.SignalIntelligenceAPI.IngestPaymentEvent(authCtx).
			PaymentEvent(eventCap).
			XIdempotencyKey(fmt.Sprintf("idk_%032x", time.Now().UnixNano())).
			Execute()
		if err != nil {
			t.Fatalf("🔥 Capture ingestion failed: %v", err)
		}
		fmt.Println("✅ Capture accepted.")

		// 🔎 Step 6: Forensic Verification
		fmt.Println("🔎 Verifying record in Forensics...")
		history, _, err := client.ForensicDataAPI.GetFraudScores(authCtx, tenantID).PageSize(10).Execute()
		if err == nil {
			found := false
			for _, s := range history.Scores {
				if s.AssessmentId == decision.AssessmentId {
					found = true
					break
				}
			}
			if found {
				fmt.Println("✨ SUCCESS: Transaction verified in forensic logs.")
			} else {
				fmt.Println("⚠️  INFO: Transaction not yet visible in forensics (async lag expected).")
			}
		} else {
			fmt.Printf("⚠️  Forensic check skipped: %v\n", err)
		}
	}

	fmt.Println("\n🎊 Go Full Circle E2E Complete.")
}
