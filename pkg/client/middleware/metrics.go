package middleware

import (
	"net/http"
	"strconv"
	"strings"
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

		if path := matchMetricPath(r); path != "" {
			m.HTTPRequestTotal.WithLabelValues(r.Method, path, status).Inc()
			m.HTTPRequestDuration.WithLabelValues(r.Method, path, status).Observe(duration)
		}

	})
}

func matchMetricPath(r *http.Request) string {
	if r.Pattern != "" {
		return r.Pattern
	}
	switch {
	case r.URL.Path == "/user/login":
		return "/user/login"
	case r.URL.Path == "/admin":
		return "/admin"
	case r.URL.Path == "/admin/review/update":
		return "/admin/review/update"
	case r.URL.Path == "/user/register":
		return "/user/register"
	case r.URL.Path == "/user/create":
		return "/user/create"
	case r.URL.Path == "/task/review":
		return "/task/review"
	case r.URL.Path == "/api/v1/grader/callback":
		return "/api/v1/grader/callback"
	case strings.HasPrefix(r.URL.Path, "/user/account/"):
		return "/user/account/{userID}"
	default:
		return ""
	}
}
