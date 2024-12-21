# HTTP Middleware

A collection of HTTP middleware components for Go web applications with a focus on logging and request tracing.

## Features

- Request ID generation and tracking
- Audit logging with configurable headers
- Structured logging using `slog`
- Easy integration with standard `http.Handler` interface
- Compatible with Chi router and other frameworks
- Fully configurable through options pattern

## Installation

```bash
go get github.com/yourusername/middleware
```

## Quick Start

```go
package main

import (
    "log/slog"
    "net/http"
    "os"

    "github.com/yourusername/middleware"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    finalHandler := middleware.RequestID()(
        middleware.Audit(
            middleware.WithHeaders([]string{"X-Request-ID"}),
            middleware.WithLogger(logger),
        )(handler),
    )

    http.ListenAndServe(":8080", finalHandler)
}
```

## Middleware Components

### RequestID Middleware

Adds a unique request ID to each incoming request.

```go
// Default usage
router.Use(middleware.RequestID())

// With custom configuration
router.Use(middleware.RequestID(
    middleware.WithHeaderName("X-Transaction-ID"),
    middleware.WithContextKey("transactionID"),
))
```

### Audit Middleware

Logs specified request headers using structured logging.

```go
router.Use(middleware.Audit(
    middleware.WithHeaders([]string{"X-Request-ID", "Consumer"}),
    middleware.WithLogger(logger),
    middleware.WithMessage("Incoming API request"),
))
```

## Examples

See the `examples` directory for more detailed usage examples:

- `examples/basic`: Basic usage with standard `http.Handler`
- `examples/customlogger`: Custom logging configuration

## License

MIT