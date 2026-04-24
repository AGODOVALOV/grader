package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/AGODOVALOV/grader/pkg/client/session"
	"github.com/AGODOVALOV/grader/pkg/logger"
)

// UploadTask godoc
// @Summary Submit a task for review
// @Description Uploads a file for a selected task and creates a new review request
// @Tags task
// @Accept multipart/form-data
// @Produce html
// @Param taskNumber formData int true "Task number"
// @Param taskFile formData file true "Source file"
// @Success 303 {string} string "See Other"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /task/review [post]
func (h *UserHandler) UploadTask(w http.ResponseWriter, r *http.Request) {
	currSession, ok := r.Context().Value(session.SessionKey).(session.Session)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	taskNumStr := r.FormValue("taskNumber")

	file, header, err := r.FormFile("taskFile")
	if err != nil {
		logger.Z(r.Context()).Error(r.Context(), "upload file", err.Error())
		http.Error(w, "request error", http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			logger.Z(r.Context()).Error(r.Context(), "upload file", err.Error())
		}
	}(file)

	taskNum, err := strconv.Atoi(taskNumStr)
	if err != nil {
		logger.Z(r.Context()).Error(r.Context(), "upload file", err.Error())
		http.Error(w, "request error", http.StatusBadRequest)
		return
	}

	objectName := fmt.Sprintf("review_%d_%d_%s", currSession.UserID, taskNum, header.Filename)

	_, err = h.Service.CreateNewReview(
		r.Context(),
		currSession.UserID,
		taskNum,
		objectName,
		file,
		header.Size)

	if err != nil {
		logger.Z(r.Context()).Error(r.Context(), "upload file", err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/user/account/%d", currSession.UserID), http.StatusFound)
}
