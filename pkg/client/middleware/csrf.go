package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/AGODOVALOV/grader/pkg/logger"
)

// /task/review → CSRF required for POST
// /user/login → CSRF required for POST
// /user/register → CSRF required for POST

var noCSRFUrls = map[string]struct{}{
	"/api/v1/grader/callback": {},
}

func CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noCSRFUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		if isSafeMethod(r.Method) {
			applyCSRFCookie(w, r)
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("csrf_token")
		if err != nil || cookie.Value == "" {
			logger.Z(r.Context()).Error(r.Context(), "csrf token missing", err.Error())
			http.Error(w, "csrf token missing", http.StatusForbidden)
			return
		}

		requestToken := r.Header.Get("X-CSRF-Token")
		if requestToken == "" {
			requestToken = r.FormValue("_csrf")
		}

		if requestToken == "" || requestToken != cookie.Value {
			logger.Z(r.Context()).Error(r.Context(), "check csrf token", "token invalid")
			http.Error(w, "csrf token invalid", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func applyCSRFCookie(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("csrf_token"); err == nil && cookie.Value != "" {
		return
	}

	tokenValue, err := genCSRFToken()
	if err != nil {
		logger.Z(r.Context()).Error(r.Context(), "failed to generate CSRF token", err.Error())
		http.Error(w, "CSRF error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    tokenValue,
		Path:     "/",
		HttpOnly: false,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((24 * time.Hour).Seconds()),
	})
}

func genCSRFToken() (string, error) {
	buf := make([]byte, 32)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}
