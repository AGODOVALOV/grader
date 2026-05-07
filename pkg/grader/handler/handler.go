package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/AGODOVALOV/grader/pkg/dto"
	"github.com/AGODOVALOV/grader/pkg/grader/usecase"
	"github.com/AGODOVALOV/grader/pkg/grader/workerpool"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/storage/s3"
	"github.com/AGODOVALOV/grader/pkg/token"
)

type GraderHandler struct {
	GraderService *usecase.GraderService
}

func NewGraderHandler(fStorage *s3.FileStorage, tknMaker token.Maker, wp *workerpool.WorkerPool) *GraderHandler {
	return &GraderHandler{
		GraderService: usecase.NewGraderService(fStorage, tknMaker, wp),
	}
}

func (h *GraderHandler) Grade(w http.ResponseWriter, r *http.Request) {
	var targetPayload dto.GraderPayload

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Z(r.Context()).Error(r.Context(), "failed to close request body", err.Error())
		}
	}(r.Body)

	err = json.Unmarshal(body, &targetPayload)

	if err != nil {
		logger.Z(r.Context()).Error(r.Context(), "failed to unmarshal request body", err.Error())
		http.Error(w, "data format error", http.StatusInternalServerError)
		return
	}

	select {
	case h.GraderService.WP.Tasks <- &targetPayload:
		w.WriteHeader(http.StatusOK)
		return
	case <-r.Context().Done():
		http.Error(w, "request cancelled", http.StatusRequestTimeout)
		return
	}
}
