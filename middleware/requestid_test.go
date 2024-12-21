// middleware/requestid_test.go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const testHeaderName = "X-Test-Request-ID"
const testContextKey = "test-request-id"

func TestRequestID_Default(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(RequestIDKey)
		if requestID == nil {
			t.Fatal("Expected request ID in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestID()
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

	if resp.Header.Get(DefaultRequestIDHeader) == "" {
		t.Fatal("Expected request ID header to be set")
	}
}

func TestRequestID_CustomHeaderName(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(RequestIDKey)
		if requestID == nil {
			t.Fatal("Expected request ID in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestID(WithHeaderName(testHeaderName))
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

	if resp.Header.Get(testHeaderName) == "" {
		t.Fatal("Expected custom request ID header to be set")
	}
}

func TestRequestID_CustomContextKey(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(testContextKey)
		if requestID == nil {
			t.Fatal("Expected request ID in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestID(WithContextKey(testContextKey))
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
