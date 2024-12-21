// pkg/middleware/types.go
package middleware

type ContextKey string

const (
	RequestIDKey ContextKey = "requestID"
	UserIDKey    ContextKey = "userID"
)
