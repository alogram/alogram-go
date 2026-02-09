// Copyright (c) 2025 Alogram Inc.
// All rights reserved.

package alogram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// WebhookVerifier provides utility methods to verify the authenticity of webhooks sent by Alogram.
type WebhookVerifier struct{}

// Verify checks the HMAC-SHA256 signature of a webhook payload.
func (v WebhookVerifier) Verify(payload []byte, headerSignature string, secret string) (bool, error) {
	if headerSignature == "" || secret == "" {
		return false, fmt.Errorf("missing signature or secret")
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// Use constant time comparison to prevent timing attacks
	if hmac.Equal([]byte(expectedSignature), []byte(headerSignature)) {
		return true, nil
	}

	return false, fmt.Errorf("invalid webhook signature")
}
