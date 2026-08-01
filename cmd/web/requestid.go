package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type contextKey string

const requestIDKey contextKey = "requestID"

func newRequestID() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}

func withRequestID(r *http.Request, id string) *http.Request {
	ctx := context.WithValue(r.Context(), requestIDKey, id)
	return r.WithContext(ctx)
}
