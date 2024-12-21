package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/vhellman/lw-router/middleware"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func main() {
	// Create mux
	mux := http.NewServeMux()

	// Health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Add debug logging to see what headers are present
		resp := HealthResponse{Status: "healthy"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// Wrap mux with middleware
	handler := middleware.Audit(
		middleware.WithHeaders([]string{"X-Correlation-ID", "User-Agent"}),
		middleware.WithMessage("API Request"),
	)(
		middleware.RequestID()(mux),
	)

	// Start server
	slog.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
