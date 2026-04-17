package handler

import "html/template"

// UserHandler handles user-related operations.
type UserHandler struct {
	template *template.Template
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(t *template.Template) UserHandler {
	return UserHandler{
		template: t,
	}
}
