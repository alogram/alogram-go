// Copyright (c) 2025 Alogram Inc.
// All rights reserved.

package alogram

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/alogram/alogram-go/internal/payrisk_v1"
)

/**
 * 🛠️ **MockRiskClient**
 * 
 * A zero-dependency mock implementation of the Alogram Risk Client.
 * Allows developers to script decisions, inject errors, and verify requests.
 */
type MockRiskClient struct {
	mu              sync.Mutex
	Calls           []map[string]interface{}
	queuedResponses []interface{}
	defaultDecision string
	defaultScore    float32
	delay           time.Duration
}

func NewMockRiskClient() *MockRiskClient {
	return &MockRiskClient{
		defaultDecision: "approve",
		defaultScore:    0.1,
	}
}

func (m *MockRiskClient) SetDefaultDecision(decision string, score float32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.defaultDecision = decision
	m.defaultScore = score
}

func (m *MockRiskClient) QueueDecision(decision string, score float32, reason string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	scoreVal := score
	id := fmt.Sprintf("mock-%d", time.Now().UnixNano())
	now := time.Now().Format(time.RFC3339)
	
	resp := &payrisk_v1.DecisionResponse{
		Decision: decision,
		DecisionAt: now,
		RiskScore: score,
		FraudScore: &payrisk_v1.FraudScore{
			Score: &scoreVal,
			RiskLevel: "low",
		},
		AssessmentId: id,
	}
	
	if reason != "" {
		resp.Reasons = []payrisk_v1.ReasonDetail{
			{
				Code: "MOCK_CODE",
				Category: "behavior",
				DisplayName: "Mock Reason",
				Description: &reason,
			},
		}
	}
	m.queuedResponses = append(m.queuedResponses, resp)
}

func (m *MockRiskClient) QueueError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queuedResponses = append(m.queuedResponses, err)
}

func (m *MockRiskClient) SetDelay(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.delay = d
}

func (m *MockRiskClient) CallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Calls)
}

func (m *MockRiskClient) handleCall(method string, request interface{}, ik string, tid string) (interface{}, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, map[string]interface{}{
		"method": method,
		"request": request,
		"ik": ik,
		"tid": tid,
	})
	delay := m.delay
	var resp interface{}
	if len(m.queuedResponses) > 0 {
		resp = m.queuedResponses[0]
		m.queuedResponses = m.queuedResponses[1:]
	}
	m.mu.Unlock()

	if delay > 0 {
		time.Sleep(delay)
	}

	if resp != nil {
		if err, ok := resp.(error); ok {
			return nil, err
		}
		return resp, nil
	}

	return nil, nil
}

func (m *MockRiskClient) CheckRisk(ctx context.Context, req payrisk_v1.CheckRequest, ik string, tid string) (*payrisk_v1.DecisionResponse, error) {
	res, err := m.handleCall("CheckRisk", req, ik, tid)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res.(*payrisk_v1.DecisionResponse), nil
	}

	id := fmt.Sprintf("mock-%d", time.Now().UnixNano())
	now := time.Now().Format(time.RFC3339)
	scoreVal := m.defaultScore
	msg := "default_mock_response"

	return &payrisk_v1.DecisionResponse{
		Decision: m.defaultDecision,
		DecisionAt: now,
		RiskScore: m.defaultScore,
		FraudScore: &payrisk_v1.FraudScore{
			Score: &scoreVal,
			RiskLevel: "low",
		},
		AssessmentId: id,
		Reasons: []payrisk_v1.ReasonDetail{
			{
				Code: "DEFAULT",
				Category: "behavior",
				DisplayName: "Default",
				Description: &msg,
			},
		},
	}, nil
}

func (m *MockRiskClient) IngestSignals(ctx context.Context, req payrisk_v1.SignalsRequest, ik string, tid string) error {
	_, err := m.handleCall("IngestSignals", req, ik, tid)
	return err
}

func (m *MockRiskClient) IngestEvent(ctx context.Context, event payrisk_v1.PaymentEvent, ik string, tid string) error {
	_, err := m.handleCall("IngestEvent", event, ik, tid)
	return err
}