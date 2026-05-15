package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AGODOVALOV/grader/pkg/client/metrics/metrics"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func CollectMetricsMiddleware(m *metrics.CustomMetrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &statusResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(ww, r)

		status := strconv.Itoa(ww.statusCode)
		duration := time.Since(start).Seconds()

		m.HTTPRequestTotal.WithLabelValues(r.Method, r.Pattern, status).Inc()
		m.HTTPRequestDuration.WithLabelValues(r.Method, r.Pattern, status).Observe(duration)
	})
}
