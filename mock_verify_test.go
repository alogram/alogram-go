// Copyright (c) 2025 Alogram Inc.
package alogram

import (
	"context"
	"testing"
	"github.com/alogram/alogram-go/internal/payrisk_v1"
)

func TestMockRiskClient(t *testing.T) {
	mock := NewMockRiskClient()
	
	// Test Default
	resp, _ := mock.CheckRisk(context.Background(), payrisk_v1.CheckRequest{}, "ik", "tid")
	if resp.Decision != "approve" {
		t.Errorf("Expected approve, got %s", resp.Decision)
	}

	// Test Queued
	mock.QueueDecision("decline", 0.99, "test_reason")
	resp2, _ := mock.CheckRisk(context.Background(), payrisk_v1.CheckRequest{}, "ik", "tid")
	if resp2.Decision != "decline" || *resp2.FraudScore.Score != 0.99 {
		t.Errorf("Expected decline 0.99, got %s %f", resp2.Decision, *resp2.FraudScore.Score)
	}

	// Test Call Count
	if mock.CallCount() != 2 {
		t.Errorf("Expected 2 calls, got %d", mock.CallCount())
	}
}
