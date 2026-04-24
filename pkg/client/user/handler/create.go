package handler

import (
	"net/http"

	"github.com/AGODOVALOV/grader/pkg/common"
)

// CreateUser godoc
// @Summary Register a new user
// @Description Creates a new user account from submitted form data
// @Tags user
// @Accept application/x-www-form-urlencoded
// @Produce html
// @Param login formData string true "Login"
// @Param name formData string true "Name"
// @Param password formData string true "Password"
// @Param confirm_password formData string true "Confirm password"
// @Success 200
// @Failure 400
// @Failure 500
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.Service.CreateUser(r.Context(), login, name, hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusFound)
}
