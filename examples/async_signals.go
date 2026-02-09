// Copyright (c) 2025 Alogram Inc.
// Example: Async Signal Ingestion

package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alogram/payrisk-go"
	"github.com/alogram/payrisk-go/internal/payrisk_v1"
)

func AsyncSignalsExample() {
	client, err := alogram.NewAlogramRiskClient(alogram.ClientOptions{
		BaseURL: "https://api.alogram.ai",
		APIKey:  "sk_test_...",
	})
	if err != nil {
		log.Fatalf("❌ Failed to create client: %v", err)
	}

	fmt.Println("👤 User user_99 performed an action.")

	go func(userId string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := payrisk_v1.SignalsRequest{
			SignalsInteractionVariant: &payrisk_v1.SignalsInteractionVariant{
				SignalType: "interaction",
				Entities: payrisk_v1.EntityIds{
					EndCustomerId: payrisk_v1.PtrString(userId),
				},
				Interactions: []payrisk_v1.Interaction{
					{
						InteractionType: payrisk_v1.INTERACTIONTYPEENUM_PAGE_VIEW,
						Timestamp:       payrisk_v1.PtrString(time.Now().Format(time.RFC3339)),
					},
				},
			},
		}

		err := client.IngestSignals(ctx, req, "", "")
		if err != nil {
			log.Printf("❌ Failed to ingest signal: %v", err)
			return
		}
		log.Printf("✅ Signal page_view ingested for %s", userId)
	}("user_99")

	time.Sleep(1 * time.Second)
}
