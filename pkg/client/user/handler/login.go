// Package handler creates client handlers
package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AGODOVALOV/grader/pkg/common"
)

// LoginUserRequest represents the payload for a user login containing login credentials.
type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginUserResponse represents the response after a successful user login containing a JWT token.
type LoginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

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
// @Router /user/login [get]
func (h *UserHandler) Login(w http.ResponseWriter, _ *http.Request) {
	err := h.template.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		fmt.Println(err)
	}
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
// @Router /user/login [post]
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	userID, err := h.Service.CheckUserLogin(r.Context(), login, password)

	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, common.ErrIncorrectPassword) {
			http.Error(w, "Incorrect password", http.StatusUnauthorized)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isAdmin, err := h.Service.CheckUserIsAdminByLogin(r.Context(), login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create token
	jwtToken, payload, err := h.Service.GetNewToken(userID, login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// use in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  payload.ExpiredAt,
	})

	if isAdmin {
		http.Redirect(w, r, "/admin", http.StatusFound)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/user/account/%d", userID), http.StatusFound)

}
