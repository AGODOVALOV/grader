// Package user provides user-related functionality.
package user //nolint:revive // package name is ok

import (
	"html/template"

	"github.com/AGODOVALOV/grader/pkg/client/user/handler"
	"github.com/AGODOVALOV/grader/pkg/client/user/usecase"
)

// User represents a user.
type User struct {
	Handler *handler.UserHandler
}

// NewUser creates a new User instance.
func NewUser(t *template.Template, s *usecase.UserService) *User {
	h := handler.NewUserHandler(t, s)
	return &User{
		Handler: h,
	}
}
