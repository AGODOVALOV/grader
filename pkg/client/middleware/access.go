// Package middleware provides HTTP middleware.
package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/AGODOVALOV/grader/pkg/client/session"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/token"
)

// Limiter defines an interface for rate-limiting logic.
// Allow reports whether a request is permitted under the configured limits.
type Limiter interface {
	Allow() bool
}

// AccessLogWithCtx wraps an HTTP handler to log request details using the provided context.
func AccessLogWithCtx(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Z(ctx).Debug(ctx, "request", "new request", map[string]string{
			"method": r.Method,
			"path":   r.URL.Path,
			"remote": r.RemoteAddr,
			"start":  start.String(),
		})
	})
}

//nolint:gochecknoglobals // static whitelist for auth middleware
var noAuthUrls = map[string]struct{}{
	"/user/login":             {},
	"/user/register":          {},
	"/user/create":            {},
	"/swagger":                {},
	"/api/v1/grader/callback": {},
	"/metrics/":               {},
}

// Auth middleware.
func Auth(tokenMaker token.Maker, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		tokenCookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenPayload, err := tokenMaker.VerifyToken(tokenCookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		sess := session.Session{
			ID:     tokenPayload.ID,
			UserID: tokenPayload.UserID,
		}

		ctx := context.WithValue(r.Context(), session.SessionKey, sess)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GlobalRateLimit(limiter Limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/task/review" {
			next.ServeHTTP(w, r)
			return
		}

		if !limiter.Allow() {
			w.Header().Set("Retry-After", strconv.Itoa(1))
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
