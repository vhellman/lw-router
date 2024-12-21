package router

import "net/http"

type Router struct {
	middlewares []func(http.Handler) http.Handler
	handler     http.Handler
}

// New creates a new middleware router
func New() *Router {
	return &Router{}
}

// Use adds middleware to the chain
func (r *Router) Use(middleware func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middleware)
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var handler http.Handler
	if r.handler != nil {
		handler = r.handler
	} else {
		handler = http.DefaultServeMux
	}

	// Chain middleware in reverse order
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	handler.ServeHTTP(w, req)
}

// Handle sets the final handler for the router
func (r *Router) Handle(handler http.Handler) {
	r.handler = handler
}

// HandleFunc sets the final handler function for the router
func (r *Router) HandleFunc(fn http.HandlerFunc) {
	r.handler = fn
}
