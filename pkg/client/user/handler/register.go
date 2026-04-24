package handler

import (
	"fmt"
	"net/http"

	"github.com/AGODOVALOV/grader/pkg/logger"
)

// Register godoc
// @Summary Registration page
// @Description Renders the user registration form
// @Tags user
// @Accept html
// @Produce html
// @Success 200 {string} string "HTML page"
// @Failure 500 {string} string "Internal server error"
// @Router /user/register [get]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	err := h.template.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		logger.Z(r.Context()).Error(r.Context(), "render register page", err.Error())
		fmt.Println(err)
	}
}
