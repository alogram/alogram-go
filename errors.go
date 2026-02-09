// Copyright (c) 2025 Alogram Inc.
// All rights reserved.

package alogram

import "fmt"

// AlogramError is the base error for all SDK errors
type AlogramError struct {
	Message string
	Status  int
	Body    string
}

func (e *AlogramError) Error() string {
	if e.Status > 0 {
		return fmt.Sprintf("%s (Status: %d)", e.Message, e.Status)
	}
	return e.Message
}

type AuthenticationError struct{ AlogramError }
type RateLimitError struct{ AlogramError }
type ValidationError struct{ AlogramError }
type InternalServerError struct{ AlogramError }
type ScopedAccessError struct{ AlogramError }

func NewAlogramError(msg string, status int, body string) error {
	base := AlogramError{Message: msg, Status: status, Body: body}
	switch status {
	case 401:
		return &AuthenticationError{base}
	case 403:
		// 403 can be either auth or scope, but for SDK level we use ScopedAccess
		return &ScopedAccessError{base}
	case 429:
		return &RateLimitError{base}
	case 400, 422:
		return &ValidationError{base}
	default:
		if status >= 500 {
			return &InternalServerError{base}
		}
		return &base
	}
}
