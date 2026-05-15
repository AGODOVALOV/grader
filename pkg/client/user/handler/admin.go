package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AGODOVALOV/grader/pkg/client/session"
)

// Admin godoc
// @Summary Admin panel page
// @Description Renders the admin panel with all user reviews
// @Tags admin
// @Accept html
// @Produce html
// @Success 200 {string} string "HTML page"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /admin [get].
func (h *UserHandler) Admin(w http.ResponseWriter, r *http.Request) {
	currSession, ok := r.Context().Value(session.SessionKey).(session.Session)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	isAdmin, err := h.Service.CheckUserIsAdminByUserID(r.Context(), currSession.UserID)
	if err != nil {
		writeHTTPError(r, w, err)
	}

	if !isAdmin {
		http.Error(w, "not admin", http.StatusUnauthorized)
		return
	}

	data, err := h.Service.GetReviews(r.Context())
	if err != nil {
		writeHTTPError(r, w, err)
		return
	}

	err = h.template.ExecuteTemplate(w, "admin.html", data)
	if err != nil {
		logErrorRequestWithDump(r, err)
		http.Error(w, ErrTemplateRender.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateReviewAdmin godoc
// @Summary Update review status
// @Description Updates review status from the admin panel and refreshes the page
// @Tags admin
// @Accept application/x-www-form-urlencoded
// @Produce html
// @Param reviewID formData int true "Review ID"
// @Param status formData string true "New review status"
// @Success 303 {string} string "See Other"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /admin/review/update [post].
func (h *UserHandler) UpdateReviewAdmin(w http.ResponseWriter, r *http.Request) {
	currSession, ok := r.Context().Value(session.SessionKey).(session.Session)
	if !ok {
		logErrorRequestWithDump(r, errors.New("session not found"))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	isAdmin, err := h.Service.CheckUserIsAdminByUserID(r.Context(), currSession.UserID)
	if err != nil {
		writeHTTPError(r, w, err)
	}

	if !isAdmin {
		logErrorRequestWithDump(r, errors.New("not admin"))
		http.Error(w, "not admin", http.StatusUnauthorized)
		return
	}

	reviewIDstr := r.FormValue("reviewID")
	status := r.FormValue("status")

	reviewID, err := strconv.Atoi(reviewIDstr)
	if err != nil {
		logErrorRequestWithDump(r, err)
		http.Error(w, "reviewID is not a number", http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateReviewStatus(r.Context(), int64(reviewID), status)
	if err != nil {
		writeHTTPError(r, w, err)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}
