package handler

import (
	"net/http"

	"github.com/AGODOVALOV/grader/pkg/client/session"
	"github.com/AGODOVALOV/grader/pkg/logger"
)

type AdminReviewData struct {
	ID        int64
	UserLogin string
	TaskTitle string
	Status    string
	Message   string
	FileName  string
	CreatedAt string
	UpdatedAt string
}

type AdminPageData struct {
	Reviews []AdminReviewData
}

func (h *UserHandler) Admin(w http.ResponseWriter, r *http.Request) {
	currSession, ok := r.Context().Value(session.SessionKey).(session.Session)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	isAdmin, err := h.Service.CheckUserIsAdminByUserID(r.Context(), currSession.UserID)
	if err != nil || !isAdmin {
		http.Error(w, "not admin", http.StatusUnauthorized)
		return
	}

	data, err := h.Service.GetReviews(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.template.ExecuteTemplate(w, "admin.html", data)
	if err != nil {
		logger.Z(r.Context()).Error(r.Context(), "render account page", err.Error())
	}
}
