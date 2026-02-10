// Copyright (c) 2025 Alogram Inc.
// Example: Webhook Verification

package examples

import (
	"fmt"
	"io"
	"net/http"

	"github.com/alogram/alogram-go"
)

const WEBHOOK_SECRET = "whsec_..."

func HandleWebhookExample(w http.ResponseWriter, r *http.Request) {
	signature := r.Header.Get("x-alogram-signature")
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Read error", http.StatusInternalServerError)
		return
	}

	verifier := alogram.WebhookVerifier{}
	isValid, err := verifier.Verify(payload, signature, WEBHOOK_SECRET)

	if !isValid || err != nil {
		fmt.Printf("🛑 Invalid Signature: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println("✅ Webhook Verified!")
	w.WriteHeader(http.StatusOK)
}
