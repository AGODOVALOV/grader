package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AGODOVALOV/grader/pkg/common"
)

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Message string
}

// Login godoc
// @Summary Login page
// @Description Renders the login form
// @Tags user
// @Accept html
// @Produce html
// @Success 200 {string} string "HTML page"
// @Failure 500 {string} string "Internal server error"
// @Router /user/login [get].
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	err := h.template.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		logErrorRequestWithDump(r, err)
	}
	return
}

// LoginUser godoc
// @Summary User login
// @Description Authenticates a user with login and password, sets access cookie, and redirects to the account page
// @Tags user
// @Accept application/x-www-form-urlencoded
// @Produce html
// @Param login formData string true "Login"
// @Param password formData string true "Password"
// @Success 303 {string} string "See Other"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /user/login [post].
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	userID, err := h.Service.CheckUserLogin(r.Context(), login, password)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrRecordNotFound):
			h.MetricsCollector.Metrics.LoginAttemptsTotal.WithLabelValues("invalid_credentials").Inc()
			http.Error(w, "User not found", http.StatusNotFound)
		case errors.Is(err, common.ErrIncorrectPassword):
			h.MetricsCollector.Metrics.LoginAttemptsTotal.WithLabelValues("invalid_credentials").Inc()
			http.Error(w, "Incorrect password", http.StatusUnauthorized)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		logErrorRequestWithDump(r, err)
		return
	}

	isAdmin, err := h.Service.CheckUserIsAdminByLogin(r.Context(), login)
	if err != nil {
		writeHTTPError(r, w, err)
		return
	}

	// create token
	jwtToken, payload, err := h.Service.GetNewToken(userID, login)
	if err != nil {
		h.MetricsCollector.Metrics.LoginAttemptsTotal.WithLabelValues("token_error").Inc()
		logErrorRequestWithDump(r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// use in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  payload.ExpiredAt,
	})

	if isAdmin {
		http.Redirect(w, r, "/admin", http.StatusFound)
		return
	}

	h.MetricsCollector.Metrics.LoginAttemptsTotal.WithLabelValues("success").Inc()

	http.Redirect(w, r, fmt.Sprintf("/user/account/%d", userID), http.StatusFound)
}
