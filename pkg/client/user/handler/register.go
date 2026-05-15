package handler

import (
	"net/http"
)

// Register godoc
// @Summary Registration page
// @Description Renders the user registration form
// @Tags user
// @Accept html
// @Produce html
// @Success 200 {string} string "HTML page"
// @Failure 500 {string} string "Internal server error"
// @Router /user/register [get].
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	err := h.template.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		logErrorRequestWithDump(r, err)
		http.Error(w, ErrTemplateRender.Error(), http.StatusInternalServerError)
	}
	return
}
