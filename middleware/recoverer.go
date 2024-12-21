// pkg/middleware/recoverer.go
package middleware

import (
	"log"
	"net/http"
)

// Recoverer recovers from panics
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Safely try to get the request ID if it exists
				requestID := "unknown"
				if id := r.Context().Value(RequestIDKey); id != nil {
					if strID, ok := id.(string); ok {
						requestID = strID
					}
				}

				log.Printf("[%s] PANIC: %v", requestID, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
