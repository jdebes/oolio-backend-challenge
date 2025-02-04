package middleware

import (
	"net/http"
	"strings"
)

const (
	// Should be placed in secrets store
	validApiSecret = "apitest"
	apiKeyHeader   = "api_key"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(apiKeyHeader)

		if apiKey == "" {
			http.Error(w, "Missing API key", http.StatusUnauthorized)
			return
		}

		if strings.TrimSpace(apiKey) != validApiSecret {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
