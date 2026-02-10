# Alogram PayRisk SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/alogram/payrisk-go.svg)](https://pkg.go.dev/github.com/alogram/payrisk-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

The official Alogram PayRisk 'Smart' SDK for Go. Engineered for high-performance financial systems that demand strict type safety, context-aware operations, and native observability.

## Features

-   **🏢 Smart Client Architecture**: Specialized `AlogramRiskClient` (Secret) and `AlogramPublicClient` (Public) types.
-   **🛡️ Automated Identity**: Automatic injection of `x-api-key`, `Authorization`, and tenant context.
-   **🔄 Built-in Resiliency**: Idiomatic retries with configurable timeouts.
-   **🕵️ OpenTelemetry Native**: Integrated tracing via `go.opentelemetry.io/otel`.
-   **🧩 Type Safe**: Strongly-typed requests and responses using standard Go idioms.

## Installation

```bash
go get github.com/alogram/payrisk-go
```

## Quick Start

### Evaluate Risk (Server-Side)

Use the `AlogramRiskClient` with your Secret Key (`sk_...`) to evaluate risk for a purchase.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alogram/payrisk-go"
	"github.com/alogram/payrisk-go/v1"
)

func main() {
	// Initialize the smart client
	client, err := alogram.NewAlogramRiskClient(alogram.ClientOptions{
		APIKey:   "sk_live_your_secret_key",
		TenantID: "tenant_123",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Build a risk check request
	req := payrisk_v1.CheckRequest{
		Purchase: &payrisk_v1.Purchase{
			Amount:   99.99,
			Currency: "USD",
		},
		Identity: &payrisk_v1.Identity{
			Email: "customer@example.com",
		},
	}

	// Perform the check with context support and automatic retries
	decision, err := client.CheckRisk(context.Background(), req, "", "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Risk Decision: %s\n", decision.Decision)
}
```

## Scoped Authentication

Alogram enforces strict key scoping at the type level:

-   **`AlogramRiskClient`**: Requires a Secret Key (`sk_...`). Returns an error if initialized with a Publishable Key.
-   **`AlogramPublicClient`**: Requires a Publishable Key (`pk_...`). Returns an error if initialized with a Secret Key.

## Documentation

For full API reference, visit [developers.alogram.ai](https://developers.alogram.ai).

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.