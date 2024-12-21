// pkg/middleware/audit.go
package middleware

/**
	ex usage:
	router.Use(middleware.Audit(
    	middleware.WithHeaders([]string{"X-Request-ID", "Consumer"}),
    	middleware.WithMessage("Incoming API request"),
	))

	// With custom logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	router.Use(middleware.Audit(
		middleware.WithHeaders([]string{"X-Request-ID", "Consumer"}),
		middleware.WithLogger(logger),
	))

*/

import (
	"log/slog"
	"net/http"
)

type auditOptions struct {
	headerNames []string
	logger      *slog.Logger
	message     string
}

type AuditOption func(*auditOptions)

// WithHeaders specifies which headers to include in the audit log
func WithHeaders(headers []string) AuditOption {
	return func(o *auditOptions) {
		o.headerNames = headers
	}
}

// WithLogger sets a custom slog.Logger instance
func WithLogger(logger *slog.Logger) AuditOption {
	return func(o *auditOptions) {
		o.logger = logger
	}
}

// WithMessage sets a custom message for the log entry
func WithMessage(msg string) AuditOption {
	return func(o *auditOptions) {
		o.message = msg
	}
}

// Audit creates a middleware that logs specified request headers
func Audit(opts ...AuditOption) func(http.Handler) http.Handler {
	options := &auditOptions{
		headerNames: []string{},
		logger:      slog.Default(),
		message:     "Request headers", // default message
	}

	for _, opt := range opts {
		opt(options)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Collect header values
			attrs := make([]any, 0, len(options.headerNames))
			for _, headerName := range options.headerNames {
				if value := r.Header.Get(headerName); value != "" {
					attrs = append(attrs, headerName, value)
				}
			}

			// Log headers with slog
			options.logger.InfoContext(r.Context(),
				options.message,
				attrs...,
			)

			next.ServeHTTP(w, r)
		})
	}
}
