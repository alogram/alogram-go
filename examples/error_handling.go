// Copyright (c) 2025 Alogram Inc.
// Example: Production Error Handling

package examples

import (
	"context"
	"fmt"
	"log"

	"github.com/alogram/alogram-go"
	"github.com/alogram/alogram-go/internal/payrisk_v1"
)

func ErrorHandlingExample(client *alogram.AlogramRiskClient, req payrisk_v1.CheckRequest) {
	decision, err := client.CheckRisk(context.Background(), req, "", "")
	if err != nil {
		switch e := err.(type) {
		case *alogram.AuthenticationError:
			log.Fatalf("🚨 AUTH FAILURE: %v", e)
		case *alogram.ValidationError:
			log.Printf("❌ INVALID REQUEST: %v", e)
			return
		case *alogram.RateLimitError, *alogram.InternalServerError:
			log.Printf("⚠️ SYSTEM DEGRADED: %v", e)
			return
		default:
			log.Printf("🔥 UNEXPECTED ERROR: %v", err)
			return
		}
	}

	fmt.Printf("Decision: %s\n", decision.Decision)
}
