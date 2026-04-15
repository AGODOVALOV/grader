package handler

import (
	"net/http"
)

type Router struct {
	Mux *http.ServeMux
}

func NewRouter() *Router {
	mux := http.NewServeMux()

	return &Router{
		Mux: mux,
	}
}
