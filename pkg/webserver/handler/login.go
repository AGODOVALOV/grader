// Package handler creates webserver handlers
package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

// LoginRequest represents the payload for user login containing login credentials.
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginResponse represents the response after successful user login containing a JWT token.
type LoginResponse struct {
	Token string `json:"token"`
}

// UserHandler handles user-related operations.
type UserHandler struct {
	template *template.Template
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(t *template.Template) UserHandler {
	return UserHandler{
		template: t,
	}
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Message string
}

// Login godoc
// @Summary User login
// @Description Login with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Router /user/login [post].
func (h *UserHandler) Login(w http.ResponseWriter, _ *http.Request) {
	err := h.template.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		fmt.Println(err)
	}
}
