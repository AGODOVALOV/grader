package user

import (
	"html/template"

	"github.com/AGODOVALOV/grader/pkg/client/user/handler"
	"github.com/AGODOVALOV/grader/pkg/client/user/usecase"
)

type User struct {
	Handler *handler.UserHandler
}

func NewUser(t *template.Template, s *usecase.UserService) *User {
	h := handler.NewUserHandler(t, s)
	return &User{
		Handler: h,
	}
}
