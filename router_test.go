package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNew tests the New function
func TestNew(t *testing.T) {
	router := New()
	if router == nil {
		t.Fatal("Expected router to be non-nil")
	}
}

// TestUse tests the Use function
func TestUse(t *testing.T) {
	router := New()
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
	router.Use(middleware)
	if len(router.middlewares) != 1 {
		t.Fatalf("Expected 1 middleware, got %d", len(router.middlewares))
	}
}

// TestServeHTTP tests the ServeHTTP function
func TestServeHTTP(t *testing.T) {
	router := New()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Handle(handler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// TestHandle tests the Handle function
func TestHandle(t *testing.T) {
	router := New()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	router.Handle(handler)
	if router.handler == nil {
		t.Fatal("Expected handler to be set")
	}
}

// TestHandleFunc tests the HandleFunc function
func TestHandleFunc(t *testing.T) {
	router := New()
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {}
	router.HandleFunc(handlerFunc)
	if router.handler == nil {
		t.Fatal("Expected handler to be set")
	}
}
