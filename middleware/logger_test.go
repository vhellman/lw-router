// middleware/logger_test.go
package middleware

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testRequestID = "test-request-id"

func TestLogger(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := Logger(handler)
	server := httptest.NewServer(middleware)
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)
	ctx := context.WithValue(req.Context(), RequestIDKey, testRequestID)
	req = req.WithContext(ctx)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(nil)
	}()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	logOutput := buf.String()
	if !contains(logOutput, "Starting") || !contains(logOutput, "Completed") {
		t.Fatal("Expected log output to contain 'Starting' and 'Completed'")
	}
}

func TestNewResponseWriter(t *testing.T) {
	recorder := httptest.NewRecorder()
	rw := newResponseWriter(recorder)

	if rw.statusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, rw.statusCode)
	}

	if rw.ResponseWriter != recorder {
		t.Fatal("Expected ResponseWriter to be set correctly")
	}
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	recorder := httptest.NewRecorder()
	rw := newResponseWriter(recorder)

	rw.WriteHeader(http.StatusNotFound)
	if rw.statusCode != http.StatusNotFound {
		t.Fatalf("Expected status code %d, got %d", http.StatusNotFound, rw.statusCode)
	}

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("Expected recorder code %d, got %d", http.StatusNotFound, recorder.Code)
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
