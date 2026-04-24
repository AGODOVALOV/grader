// Package handler creates client handlers
package handler

import (
	"net/http"

	"github.com/AGODOVALOV/grader/pkg/client/session"
	"github.com/AGODOVALOV/grader/pkg/logger"
)

// Account godoc
// @Summary User account page
// @Description Renders the personal account page for the authenticated user
// @Tags user
// @Accept html
// @Produce html
// @Param userID path int true "User ID"
// @Success 200 {string} string "HTML page"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /user/account/{userID} [get].
func (h *UserHandler) Account(w http.ResponseWriter, r *http.Request) {
	currSession, ok := r.Context().Value(session.SessionKey).(session.Session)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	data, err := h.Service.GetReviewsByUserID(r.Context(), currSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = h.template.ExecuteTemplate(w, "account.html", data)
	if err != nil {
		logger.Z(r.Context()).Error(r.Context(), "render account page", err.Error())
		return
	}
}
