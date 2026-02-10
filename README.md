# Alogram PayRisk SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/alogram/alogram-go.svg)](https://pkg.go.dev/github.com/alogram/alogram-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

The official Alogram PayRisk 'Smart' SDK for Go. Engineered for high-throughput financial systems requiring context-aware resiliency, native observability, and ergonomic risk intelligence.

## 🚀 Features

-   **🏢 Smart Client Architecture**: Specialized clients for server-side (`AlogramRiskClient`) and edge (`AlogramPublicClient`).
-   **🛡️ Automated Identity**: Thread-safe injection of `x-api-key`, `Authorization`, and tenant headers.
-   **🔄 Built-in Resiliency**: Transparent exponential backoff and jittered retries (3 retries on 429/5xx).
-   **🕵️ OpenTelemetry Native**: Built-in tracing support for `context.Context` propagated risk decisions.
-   **🏗️ Go Idiomatic**: Uses standard `context.Context` for cancellation and timeouts.

## 📦 Installation

```bash
go get github.com/alogram/alogram-go
```

## 🛠️ Quick Start

### Evaluate Risk (Server-Side)

```go
import (
    "context"
    "fmt"
    "github.com/alogram/alogram-go"
    "github.com/alogram/alogram-go/internal/payrisk_v1"
)

func main() {
    client, _ := alogram.NewAlogramRiskClient(alogram.ClientOptions{
        APIKey: "sk_live_your_secret_key",
    })

    req := payrisk_v1.CheckRequest{
        Purchase: payrisk_v1.Purchase{Amount: 99.99, Currency: "USD"},
    }

    decision, err := client.CheckRisk(context.Background(), req, "idk_123", "trc_123")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Decision: %s\n", decision.Decision)
}
```

---

## 🛡️ Error Handling

Use `errors.As` to handle specific API error states:

```go
import "errors"

decision, err := client.CheckRisk(...)
if err != nil {
    var rateLimitErr *alogram.RateLimitError
    if errors.As(err, &rateLimitErr) {
        // 🔄 Backoff is handled by SDK, but you can log here
    }
    
    var validationErr *alogram.ValidationError
    if errors.As(err, &validationErr) {
        fmt.Println("Invalid input:", validationErr.Body)
    }
}
```

## 🔑 Authentication & Keys

Obtain your keys from the [Alogram Dashboard](https://dashboard.alogram.ai/keys).

| Key Type | Prefix | Client | Use Case |
| :--- | :--- | :--- | :--- |
| **Secret Key** | `sk_live_...` | `AlogramRiskClient` | Server-side risk checks. **Never expose.** |
| **Publishable Key** | `pk_live_...` | `AlogramPublicClient` | Frontend signals. |

## 🤖 For AI Agents (Integration Context)

If you are using an AI agent (like ChatGPT, Claude, or Gemini), provide this context:

> "I am integrating the Alogram PayRisk Go SDK. I need to call `CheckRisk` in my transaction handler. Use the `AlogramRiskClient`, handle `RateLimitError` gracefully, pass the request `context.Context` for tracing, and ensure you use the `idempotencyKey` parameter."

---

## 🧪 Local Testing & Mocking

Use the thread-safe `MockRiskClient` for unit tests:

```go
mock := alogram.NewMockRiskClient()
mock.QueueDecision("decline", 0.99, "high_risk_ip")

decision, _ := myApp.Process(mock)
fmt.Println(decision.Decision) // "decline"
```

## 🏗️ Environment Testing

### Alogram Sandbox
For safe integration testing, point your client to the Sandbox environment:
```go
client, _ := alogram.NewAlogramRiskClient(alogram.ClientOptions{
    APIKey:  "sk_test_...",
    BaseURL: "https://api-sandbox.alogram.ai",
})
```

### Local Emulator
For hermetic local testing, run the **Alogram Local Emulator**:
```bash
docker run -p 8080:8080 alogram/payrisk-emulator
```
Point your client to the local instance:
```go
client, _ := alogram.NewAlogramRiskClient(alogram.ClientOptions{
    BaseURL: "http://localhost:8080",
    APIKey:  "test",
})
```

---

## 📚 Documentation

For the full API reference, visit [docs.alogram.ai](https://docs.alogram.ai).

## ⚖️ License

Apache License 2.0. See [LICENSE](LICENSE) for details.