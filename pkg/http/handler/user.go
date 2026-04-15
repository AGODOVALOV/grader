package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

type UserHandler struct {
	template *template.Template
}

func NewUserHandler(t *template.Template) UserHandler {
	return UserHandler{
		template: t,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, _ *http.Request) {

	err := h.template.ExecuteTemplate(w, "login.html", nil)

	if err != nil {
		fmt.Println(err)
	}

}
