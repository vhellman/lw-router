// pkg/middleware/logger.go
package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger logs request details
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := r.Context().Value(RequestIDKey).(string)
		rw := newResponseWriter(w)

		log.Printf("[%s] Starting %s %s", requestID, r.Method, r.URL.Path)
		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf("[%s] Completed %s %s [%d] in %v",
			requestID, r.Method, r.URL.Path, rw.statusCode, duration,
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
