package handler

import (
	"net/http"
)

// Router represents an HTTP request router with an associated ServeMux for handling routes and requests.
type Router struct {
	Mux *http.ServeMux
}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	mux := http.NewServeMux()

	return &Router{
		Mux: mux,
	}
}
