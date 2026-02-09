# Alogram Payrisk Go SDK

The official Go client for the **Alogram Payments Risk API**. This SDK provides a robust, "smart" interface for checking fraud risk, ingesting behavioral signals, and managing payment lifecycle events.

**Key Features:**
*   **Resilient:** Built-in retries with exponential backoff for transient errors.
*   **Traceable:** Automatic injection of `x-trace-id` and `x-idempotency-key` for every request.
*   **Observable:** First-class support for **OpenTelemetry** spans and attributes.
*   **Type-Safe:** Fully typed request/response models.
*   **Secure:** Built-in webhook signature verification.

---

## 🏗️ Installation

```bash
go get github.com/alogram/payrisk-go
```

---

## 🚀 Quickstart

### 1. Initialize the Client

```go
import "github.com/alogram/payrisk-go"

client := alogram.NewAlogramRiskClient(alogram.ClientOptions{
    BaseURL: "https://api.alogram.ai",
    APIKey:  "sk_live_...",
    TenantID: "your_tenant_id",
})
```

### 2. Check Risk

```go
package main

import (
    "context"
    "fmt"
    "github.com/alogram/payrisk-go"
    "github.com/alogram/payrisk-go/internal/payrisk_v1"
)

func main() {
    client := alogram.NewAlogramRiskClient(alogram.ClientOptions{
        BaseURL: "https://api.alogram.ai",
        APIKey:  "sk_live_...",
    })

    req := payrisk_v1.CheckRequest{
        EventType: "purchase",
        Entities: &payrisk_v1.EntityIds{
            TenantId: "tid_123",
            ClientId: "cid_abc",
        },
        Purchase: &payrisk_v1.Purchase{
            Amount:   99.00,
            Currency: "USD",
        },
    }

    decision, err := client.CheckRisk(context.Background(), req, "", "")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Decision: %s | Score: %f\n", *decision.Decision, *decision.RiskScore)
}
```

---

## 📊 Observability (OpenTelemetry)

The SDK uses the standard OpenTelemetry Go API. It automatically detects and uses any configured Global Tracer.

**Captured Attributes:**
*   `alogram.idempotency_key`
*   `alogram.trace_id`
*   `alogram.decision`

---

## 🛡️ Webhook Security

Verify incoming webhooks using the built-in `WebhookVerifier`.

```go
verifier := alogram.WebhookVerifier{}
isValid, err := verifier.Verify(payloadBytes, signatureHeader, webhookSecret)
```

---

## ⚠️ Error Handling

| Exception | Description |
| :--- | :--- |
| `AuthenticationError` | Invalid API Key or Permissions. |
| `ValidationError` | Invalid request body or missing fields. |
| `RateLimitError` | Too many requests. **Automatically Retried.** |
| `InternalServerError` | Server-side issues. **Automatically Retried.** |

---

## 📦 License

Apache 2.0
