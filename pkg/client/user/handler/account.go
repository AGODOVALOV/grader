package handler

import (
	"net/http"

	"github.com/AGODOVALOV/grader/pkg/client/session"
	"github.com/AGODOVALOV/grader/pkg/logger"
)

type AccountPageData struct {
	ID    int
	Login string
	Name  string
	Tasks []TaskData
}

type TaskData struct {
	ID        int
	Title     string
	Status    string
	Message   string
	UpdatedAt string
}

func (h *UserHandler) Account(w http.ResponseWriter, r *http.Request) {
	currSession, ok := r.Context().Value(session.SessionKey).(session.Session)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	data, err := h.Service.GetReviewsByUserID(r.Context(), currSession.UserID)

	err = h.template.ExecuteTemplate(w, "account.html", data)
	if err != nil {
		logger.Z(r.Context()).Error(r.Context(), "render account page", err.Error())
	}
}
