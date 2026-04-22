package handler

import (
	"net/http"

	"github.com/AGODOVALOV/grader/pkg/common"
)

type CreateUserRequest struct {
	Login           string `json:"login"`
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

// ErrorResponse represents an error response.
type ErrorResponseCreateUser struct {
	Message string
}

// CreateUser godoc
// @Summary User login create
// @Description CreateUser with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "login request"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponseCreateUser
// @Router /user/create [post].
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	login := r.FormValue("login")
	name := r.FormValue("name")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if password != confirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	hashedPassword, err := common.HashPassword(password)

	err = h.Service.CreateUser(r.Context(), login, name, hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusFound)
}
