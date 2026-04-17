package handler

import (
	"fmt"
	"net/http"
)

func (h *UserHandler) Register(w http.ResponseWriter, _ *http.Request) {
	err := h.template.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		fmt.Println(err)
	}
}
