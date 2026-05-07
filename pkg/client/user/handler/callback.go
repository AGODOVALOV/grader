package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/AGODOVALOV/grader/pkg/dto"
	"github.com/AGODOVALOV/grader/pkg/logger"
)

const bearerPrefix = "Bearer "

// CallBack godoc
// @Summary Callback from grader
// @Description Callback from grader to process the review request with result
// @Tags callback
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dto.GraderPayloadCallback true "Grader callback payload"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Invalid callback payload"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/grader/callback [post]
func (h *UserHandler) CallBack(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	rawToken := strings.TrimPrefix(authHeader, bearerPrefix)

	_, err := h.Service.VerifyTokenCallBackToken(rawToken)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Z(r.Context()).Error(r.Context(), "failed to close request body", err.Error())
		}
	}(r.Body)

	var callbackPayload dto.GraderPayloadCallback
	if err := json.NewDecoder(r.Body).Decode(&callbackPayload); err != nil {
		http.Error(w, "invalid callback payload", http.StatusBadRequest)
		return
	}

	err = h.Service.ProcessGraderCallback(r.Context(), &callbackPayload)
	if err != nil {
		http.Error(w, "failed to process callback", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
