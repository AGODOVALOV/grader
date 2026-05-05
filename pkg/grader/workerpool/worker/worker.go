package worker

import (
	"context"
	"errors"

	"github.com/AGODOVALOV/grader/pkg/dto"
	"github.com/AGODOVALOV/grader/pkg/storage/s3"
)

type Worker struct {
	fStorage *s3.FileStorage
}

func NewWorker(fStorage *s3.FileStorage) *Worker {
	return &Worker{
		fStorage: fStorage,
	}
}

func (w *Worker) DoJob(ctx context.Context, payload *dto.GraderPayload) error {
	if len(payload.FileIDs) == 0 {
		return errors.New("no files provided")
	}

	// get file from S3
	fileData, err := w.fStorage.DownloadFile(ctx, payload.FileIDs[0].FileName)
	if err != nil {
		return err
	}

	// start docker flow for tests files
	_ = fileData

	return nil
}
