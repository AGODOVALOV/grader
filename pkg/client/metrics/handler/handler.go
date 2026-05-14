package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	Metrics http.Handler
}

func NewHandler(reg *prometheus.Registry) *Handler {
	return &Handler{
		Metrics: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
	}
}
