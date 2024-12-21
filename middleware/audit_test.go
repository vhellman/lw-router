// middleware/audit_test.go
package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"
)

func TestWithHeaders(t *testing.T) {
	// Create a buffer to capture logs if you want to inspect them
	var buf bytes.Buffer

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create logger with a proper writer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	middleware := Audit(
		WithHeaders([]string{"X-Test-Header"}),
		WithLogger(logger),
	)

	server := httptest.NewServer(middleware(handler))
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)
	req.Header.Set("X-Test-Header", "test-value")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
func TestWithLogger(t *testing.T) {
	// Create a buffer to capture logs if you want to inspect them
	var buf bytes.Buffer

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create logger with a proper writer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	middleware := Audit(
		WithLogger(logger),
	)

	server := httptest.NewServer(middleware(handler))
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
func TestWithMessage(t *testing.T) {
	// Create a buffer to capture logs if you want to inspect them
	var buf bytes.Buffer

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create logger with a proper writer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	middleware := Audit(
		WithMessage("Custom log message"),
		WithLogger(logger),
	)

	server := httptest.NewServer(middleware(handler))
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check if the custom message is in the logs
	if !bytes.Contains(buf.Bytes(), []byte("Custom log message")) {
		t.Fatalf("Expected log message to contain 'Custom log message', got %s", buf.String())
	}
}
func TestAuditWithNoOptions(t *testing.T) {
	// Create a buffer to capture logs if you want to inspect them
	var buf bytes.Buffer

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create logger with a proper writer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	middleware := Audit(
		WithLogger(logger),
	)

	server := httptest.NewServer(middleware(handler))
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check if the default message is in the logs
	if !bytes.Contains(buf.Bytes(), []byte("Request headers")) {
		t.Fatalf("Expected log message to contain 'Request headers', got %s", buf.String())
	}
}

func TestAuditWithMultipleHeaders(t *testing.T) {
	// Create a buffer to capture logs if you want to inspect them
	var buf bytes.Buffer

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create logger with a proper writer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	middleware := Audit(
		WithHeaders([]string{"X-Test-Header-1", "X-Test-Header-2"}),
		WithLogger(logger),
	)

	server := httptest.NewServer(middleware(handler))
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)
	req.Header.Set("X-Test-Header-1", "value1")
	req.Header.Set("X-Test-Header-2", "value2")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check if the headers are in the logs
	if !bytes.Contains(buf.Bytes(), []byte("X-Test-Header-1")) || !bytes.Contains(buf.Bytes(), []byte("value1")) {
		t.Fatalf("Expected log to contain 'X-Test-Header-1: value1', got %s", buf.String())
	}
	if !bytes.Contains(buf.Bytes(), []byte("X-Test-Header-2")) || !bytes.Contains(buf.Bytes(), []byte("value2")) {
		t.Fatalf("Expected log to contain 'X-Test-Header-2: value2', got %s", buf.String())
	}
}
