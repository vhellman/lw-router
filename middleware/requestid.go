// pkg/middleware/requestid.go
package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"time"
)

const DefaultRequestIDHeader = "X-Request-ID"

type requestIDOptions struct {
	headerName string
	contextKey ContextKey
}

type RequestIDOption func(*requestIDOptions)

// WithHeaderName sets a custom header name for the request ID
func WithHeaderName(name string) RequestIDOption {
	return func(o *requestIDOptions) {
		o.headerName = name
	}
}

// WithContextKey sets a custom context key for the request ID
func WithContextKey(key ContextKey) RequestIDOption {
	return func(o *requestIDOptions) {
		o.contextKey = key
	}
}

// generateUUID generates a UUID using crypto/rand
func generateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	// Set version (4) and variant (2) bits according to RFC 4122
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
}

// RequestID adds a unique request ID to the context and response headers
func RequestID(opts ...RequestIDOption) func(http.Handler) http.Handler {
	options := &requestIDOptions{
		headerName: DefaultRequestIDHeader,
		contextKey: RequestIDKey,
	}

	for _, opt := range opts {
		opt(options)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if request already has an ID
			requestID := r.Header.Get(options.headerName)
			if requestID == "" {
				var err error
				requestID, err = generateUUID()
				if err != nil {
					// If UUID generation fails, use timestamp or another fallback
					requestID = fmt.Sprintf("fallback-%d", time.Now().UnixNano())
				}
			}

			// Always set the header in both request and response
			r.Header.Set(options.headerName, requestID)
			w.Header().Set(options.headerName, requestID)

			// Add to context
			ctx := context.WithValue(r.Context(), options.contextKey, requestID)

			// Continue with the modified request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
