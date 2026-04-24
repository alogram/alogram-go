<p align="center">
  <img src="https://raw.githubusercontent.com/alogram/alogram-python/main/.github/assets/logo.png" width="200" alt="Alogram PayRisk Logo">
</p>

# Alogram PayRisk SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/alogram/alogram-go.svg)](https://pkg.go.dev/github.com/alogram/alogram-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

The official Go client for the **Alogram PayRisk Engine**. 

Alogram PayRisk is a decision management and risk orchestration engine for global commerce. It fuses machine learning, behavioral analytics, and deterministic business rules into a high-fidelity scoring pipeline designed for enterprise scale and auditability.
## 🧠 The Three-Expert Architecture

The SDK provides unified access to three specialized risk experts:

-   **Risk Scoring**: Real-time assessment and decision orchestration for purchases.
-   **Signal Intelligence**: Ingestion of behavioral telemetry and payment lifecycle events.
-   **Forensic Data**: Deep visibility into historical assessments and decision transparency.

---

## 🔐 Security: Trust Boundaries

Alogram enforces a strict separation between client-side telemetry and server-side decisioning.

| Client Type | Key Prefix | Environment | Capabilities |
| :--- | :--- | :--- | :--- |
| **`AlogramPublicClient`** | `pk_...` | Browser / Mobile | **Ingestion only.** Restricted to behavioral signals. |
| **`AlogramRiskClient`** | `sk_...` | Secure Backend | **Full access.** Authorized for risk decisions and forensics. |

> [!WARNING]
> **Never** use a Secret Key (`sk_...`) in a client-side environment. This will expose your tenant's sensitive forensic data and violate Alogram's security mandates.

---

## 🔄 Full Lifecycle Integration

Go's concurrency model makes it ideal for background signal ingestion and synchronous risk checks.

```go
import (
    "github.com/alogram/alogram-go/payrisk"
)

// 1. Initialize the Secret Client (Backend Only)
client := payrisk.NewRiskClient("sk_live_...", "tid_mycorp")

// 2. Assessment: Call before charging the card
decision, err := client.CheckRisk(ctx, payrisk.CheckRequest{...})
if err != nil || decision.Decision != "approve" {
    return
}

// 3. Lifecycle: Send the outcome back to Alogram
err = client.IngestEvent(ctx, payrisk.PaymentEvent{
    PaymentIntentId: decision.PaymentIntentId,
    EventType:       payrisk.AUTHORIZATION,
    Outcome:         payrisk.PaymentOutcome{Approved: true},
})
```

---

## 🚀 High-Performance Integration
-   **🏢 Smart Client Architecture**: Specialized clients for server-side (`AlogramRiskClient`) and edge (`AlogramPublicClient`).
-   **🛡️ Automated Identity**: Thread-safe injection of `x-api-key`, `Authorization`, and tenant headers.
-   **🔄 Built-in Resiliency**: Automatic exponential backoff and jittered retries (3 retries on 429/5xx).
-   **🕵️ Native Observability**: Built-in support for `context.Context` propagated OpenTelemetry tracing.
-   **🏗️ Go Idiomatic**: Designed for high-concurrency systems using standard Go patterns.

## 📦 Installation

```bash
go get github.com/alogram/alogram-go
```

## 🛠️ Quick Start

### Evaluate Risk (Risk Scoring Expert)

Assess a purchase in real-time. This invokes the authoritative scoring pipeline.

```go
import (
    "context"
    "fmt"
    "github.com/alogram/alogram-go"
    "github.com/alogram/alogram-go/internal/payrisk_v1"
)

func main() {
    client, _ := alogram.NewAlogramRiskClient(alogram.ClientOptions{
        APIKey: "sk_live_...",
    })

    req := payrisk_v1.CheckRequest{
        Purchase: payrisk_v1.Purchase{Amount: 99.99, Currency: "USD"},
    }

    // Perform the check with automatic tracing and retries
    decision, err := client.CheckRisk(context.Background(), req)
    if err != nil {
        panic(err)
    }
    
    // Access the authoritative blended score
    fmt.Printf("Decision: %s | Score: %.2f
", decision.Decision, decision.DecisionScore)
}
```

---

## 🚀 High-Performance Integration

To ensure sub-second risk assessment latencies and handle high-volume signal telemetry efficiently, please adhere to these network best practices:

-   **Persistent Client:** Reuse the `Client` object for all operations. 
    -   *Why:* Go's `http.Client` manages a persistent connection pool. Creating a new client for every request bypasses this pool, forcing fresh TCP/TLS handshakes.
-   **Native HTTP/2 Multiplexing:** Go natively supports HTTP/2 for `https://` endpoints. Reusing the client allows thousands of concurrent operations to be multiplexed over a single connection with minimal overhead.

## 🛡️ Error Handling & Resiliency

The SDK distinguishes between transient network issues and validation errors.

```go
import "errors"

decision, err := client.CheckRisk(...)
if err != nil {
    var validationErr *alogram.ValidationError
    if errors.As(err, &validationErr) {
        // ❌ Handle invalid input (check validationErr.Body)
    }
    
    var alogramErr *alogram.AlogramError
    if errors.As(err, &alogramErr) {
        // 🚨 Handle generic API or Auth errors
    }
}
```

## 🕵️ Observability (OpenTelemetry)

The SDK automatically integrates with OpenTelemetry if active in your `context.Context`.

```go
// Spans will automatically include 'alogram.decision' and 'alogram.trace_id'
ctx, span := tracer.Start(ctx, "checkout")
defer span.End()

decision, err := client.CheckRisk(ctx, request)
```

## 🏗️ Environments

| Environment | Base URL | Key Type |
| :--- | :--- | :--- |
| **Production** | `https://api.alogram.ai` | `sk_live_...` |
| **Sandbox** | `https://api-sandbox.alogram.ai` | `sk_test_...` |
| **Local** | `http://localhost:8080` | `test` |

---

## 🤖 For AI Agents

If using an AI agent (ChatGPT, Claude, Gemini) for integration, provide this context:
> "Use the Alogram PayRisk Go SDK. Always prefer the `AlogramRiskClient` for backend operations. Pass the `context.Context` to all API calls for proper tracing. Access the authoritative blended score via the `DecisionScore` field of the response."

---

## ⚖️ License

Apache License 2.0. See [LICENSE](LICENSE) for details.
