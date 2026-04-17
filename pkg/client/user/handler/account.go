package handler

import (
	"fmt"
	"net/http"
	"time"
)

type AccountPageData struct {
	ID    int
	Login string
	Email string
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

	userID := r.PathValue("userID")

	_ = userID

	data := AccountPageData{
		ID:    1,
		Login: "john",
		Email: "john@example.com",
		Tasks: []TaskData{
			{
				ID:        1,
				Title:     "First task",
				Status:    "done",
				Message:   "Task completed successfully",
				UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
			},
			{
				ID:        2,
				Title:     "Second task",
				Status:    "In process",
				Message:   "Task in progress",
				UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
			},
			{
				ID:        3,
				Title:     "Third task",
				Status:    "Error",
				Message:   "Task failed",
				UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
			},
		},
	}

	err := h.template.ExecuteTemplate(w, "account.html", data)
	if err != nil {
		fmt.Println(err)
	}
}
