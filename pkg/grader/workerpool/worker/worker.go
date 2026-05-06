package worker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/AGODOVALOV/grader/pkg/dto"
	"github.com/AGODOVALOV/grader/pkg/logger"
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
	var (
		vPathHost        string
		containerWorkdir string
		volumePath       string
		err              error
	)
	if len(payload.FileIDs) == 0 {
		return errors.New("no files provided")
	}

	fName := payload.FileIDs[0].FileName

	// get file from S3
	fileData, err := w.fStorage.DownloadFile(ctx, fName)
	if err != nil {
		return err
	}

	fName = fName + ".file"

	//save files to local submission storage
	switch payload.TaskID {
	case "1":
		vPathHost = filepath.Join("./infra/grader/submission/1/", fName)
		containerWorkdir = "/app/1/game"
		err = os.WriteFile(vPathHost, fileData, 0644)
		if err != nil {
			logger.Z(ctx).Error(ctx, "write file", err.Error())
			return err
		}
	case "2":
		vPathHost = filepath.Join("./infra/grader/submission/2/", fName)
		containerWorkdir = "/app/2/client"
		err = os.WriteFile(vPathHost, fileData, 0644)
		if err != nil {
			logger.Z(ctx).Error(ctx, "write file", err.Error())
			return err
		}
	}

	// start docker flow for tests files
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	//docker run --rm --network none -v "./infra/grader/submission/1/review_9_1_main.go.file:/app/1/game/main.go" --workdir /app/1/game grader:latest go test
	//docker run --rm --network none -v "./infra/grader/submission/2/review_9_2_main.go.file:/app/2/client/main.go" --workdir /app/2/client grader:latest go test

	volumePath = "./" + vPathHost + ":" + containerWorkdir + "/main.go"

	cmd := exec.CommandContext(ctx, "docker", "run",
		"--rm",
		"--network", "none",
		"--memory", "256m",
		"--cpus", "1",
		"-v", volumePath,
		"--workdir", containerWorkdir,
		"grader:latest",
		"go", "test",
	)

	out, err := cmd.CombinedOutput()
	switch {
	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		logger.Z(ctx).Error(ctx, "docker runtime deadline", "timeout")
		return ctx.Err()
	case err != nil:
		logger.Z(ctx).Error(ctx, "docker runtime", err.Error(), map[string]string{
			"output": string(out),
		})
		return err
	default:
		logger.Z(ctx).Info(ctx, "docker run", string(out))
	}

	logger.Z(ctx).Info(ctx, "docker run result", fmt.Sprintf("test for %s", fName), map[string]string{
		"result": string(out),
	})

	return nil
}
