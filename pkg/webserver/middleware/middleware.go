// Package middleware provides HTTP middleware.
package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/AGODOVALOV/grader/pkg/logger"
)

// AccessLogWithCtx wraps an HTTP handler to log request details using the provided context.
func AccessLogWithCtx(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// TODO create request scope request
		next.ServeHTTP(w, r.WithContext(ctx))
		logger.Z(ctx).Debug(ctx, "request", "new request", map[string]string{
			"method": r.Method,
			"path":   r.URL.Path,
			"remote": r.RemoteAddr,
			"start":  start.String(),
		})
	})
}
