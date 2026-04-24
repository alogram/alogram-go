package payrisk_v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Smoke_HealthCheck(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/health" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	cfg := NewConfiguration()
	cfg.Servers = ServerConfigurations{{URL: ts.URL, Description: "test"}}

	client := NewAPIClient(cfg)
	ctx := context.Background()

	resp, err := client.SystemAPI.HealthCheck(ctx).Execute()
	if err != nil {
		t.Fatalf("HealthCheck error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}
}
