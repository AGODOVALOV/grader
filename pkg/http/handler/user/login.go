package user

import (
	"fmt"
	"html/template"
	"net/http"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UserHandler struct {
	template *template.Template
}

func NewUserHandler(t *template.Template) UserHandler {
	return UserHandler{
		template: t,
	}
}

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
// @Router /user/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, _ *http.Request) {

	err := h.template.ExecuteTemplate(w, "login.html", nil)

	if err != nil {
		fmt.Println(err)
	}

}
