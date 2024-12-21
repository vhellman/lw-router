// examples/customlogger/main.go
package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/vhellman/lw-router/middleware"
)

type StatusResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

func main() {
	// Create a custom JSON logger with specific options
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				return slog.String(a.Key, "[AUDIT] "+a.Value.String())
			}
			return a
		},
	})
	logger := slog.New(logHandler)

	// Create mux
	mux := http.NewServeMux()

	// Status endpoint
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		resp := StatusResponse{
			Status:    "operational",
			Timestamp: time.Now(),
			Version:   "1.0.0",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// Define middleware
	audit := middleware.Audit(
		middleware.WithHeaders([]string{
			"X-Request-ID",
			"X-Correlation-ID",
			"Authorization",
			"User-Agent",
		}),
		middleware.WithLogger(logger),
		middleware.WithMessage("Service call"),
	)

	requestID := middleware.RequestID(
		middleware.WithHeaderName("X-Request-ID"),
	)

	// Build the chain from inside out
	// The audit middleware should be outside (executed first) to capture all headers
	// including those set by inner middleware
	handler := audit(mux)        // Inner: Audit logs the headers
	handler = requestID(handler) // Outer: RequestID sets the header

	// Start server
	logger.Info("Starting server with custom logging on :8080")
	logger.Info("Test with: curl -H 'X-Correlation-ID: test-correlation' http://localhost:8080/status")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
