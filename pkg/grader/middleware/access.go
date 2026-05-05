// Package middleware provides HTTP middleware.
package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/token"
)

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

// Auth middleware.
func Auth(tokenMaker token.Maker, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//tokenCookie, err := r.Cookie("access_token")
		//if err != nil {
		//	http.Error(w, "unauthorized", http.StatusUnauthorized)
		//	return
		//}
		//
		//tokenPayload, err := tokenMaker.VerifyToken(tokenCookie.Value)
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusUnauthorized)
		//	return
		//}
		//
		//sess := session.Session{
		//	ID:     tokenPayload.ID,
		//	UserID: tokenPayload.UserID,
		//}
		//ctx := context.WithValue(r.Context(), session.SessionKey, sess)
		//next.ServeHTTP(w, r.WithContext(ctx))

		next.ServeHTTP(w, r)
	})
}
