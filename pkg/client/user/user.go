package user

import (
	"html/template"

	"github.com/AGODOVALOV/grader/pkg/client/user/handler"
	"github.com/AGODOVALOV/grader/pkg/client/user/usecase"
	"github.com/google/uuid"
)

type User struct {
	Handler *handler.UserHandler
}

type Session struct {
	ID     uuid.UUID
	UserID int64
}

func NewUser(t *template.Template, s *usecase.UserService) *User {
	h := handler.NewUserHandler(t, s)
	return &User{
		Handler: h,
	}
}
