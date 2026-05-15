package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/AGODOVALOV/grader/pkg/dto"
	"github.com/AGODOVALOV/grader/pkg/grader/client"
	"github.com/AGODOVALOV/grader/pkg/grader/metrics"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/storage/s3"
)

type Worker struct {
	fStorage         *s3.FileStorage
	callBackClient   *client.Client
	metricsCollector *metrics.Collector
}

func NewWorker(fStorage *s3.FileStorage, callBackClient *client.Client, metricsCollector *metrics.Collector) *Worker {
	return &Worker{
		fStorage:         fStorage,
		callBackClient:   callBackClient,
		metricsCollector: metricsCollector,
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

	// avoid path traversal
	fName := filepath.Base(payload.FileIDs[0].FileName)

	// get file from S3
	fileData, err := w.fStorage.DownloadFile(ctx, fName)
	if err != nil {
		return err
	}

	w.metricsCollector.Metrics.S3DownloadsTotal.WithLabelValues("s3 download").Inc()

	fName = fName + ".tst"

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
	default:
		return errors.New("unknown task id")
	}

	// delete local after docker run
	defer func() {
		err := os.Remove(vPathHost)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			logger.Z(ctx).Error(ctx, "delete file", err.Error(), map[string]string{
				"path":   vPathHost,
				"user":   payload.UserID,
				"task":   payload.TaskID,
				"review": payload.ReviewID,
				"event":  payload.EventID,
			})
		}
	}()

	// start docker flow for tests files
	runCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// docker run --rm --network none -v "./infra/grader/submission/1/review_9_1_main.go.file:/app/1/game/main.go" --workdir /app/1/game grader:latest go test
	// docker run --rm --network none -v "./infra/grader/submission/2/review_9_2_main.go.file:/app/2/client/main.go" --workdir /app/2/client grader:latest go test

	volumePath = "./" + vPathHost + ":" + containerWorkdir + "/main.go"

	dockerStart := time.Now()

	w.metricsCollector.Metrics.DockerRunsTotal.WithLabelValues("docker run").Inc()

	cmd := exec.CommandContext(runCtx, "docker", "run",
		"--rm",
		"--network", "none",
		"--memory", "256m",
		"--cpus", "1",
		"-v", volumePath,
		"--workdir", containerWorkdir,
		"grader:latest",
		"go", "test",
	)

	w.metricsCollector.Metrics.DockerRunDuration.WithLabelValues(payload.TaskID, "docker run").Observe(time.Since(dockerStart).Seconds())

	pass := true
	out, err := cmd.CombinedOutput()
	if err != nil {
		pass = false
		logger.Z(ctx).Error(ctx, "docker runtime", err.Error(), map[string]string{
			"user":   payload.UserID,
			"task":   payload.TaskID,
			"review": payload.ReviewID,
			"event":  payload.EventID,
			"output": string(out),
		})
	}
	callBackPayload := dto.GraderPayloadCallback{
		UserID:        payload.UserID,
		TaskID:        payload.TaskID,
		ReviewID:      payload.ReviewID,
		EventID:       payload.EventID,
		Passed:        pass,
		OutputMessage: string(out),
		ErrorText:     getErrText(err),
	}

	callBackPayloadBytes, err := json.Marshal(callBackPayload)
	if err != nil {
		logger.Z(ctx).Error(ctx, "json marshal", err.Error())
		return err
	}

	logger.Z(ctx).Debug(ctx, "docker run result", fmt.Sprintf("test for %s", fName), map[string]string{
		"result": string(out),
	})

	err = w.callBackClient.DoCallbackRequestWithRetry(ctx, callBackPayloadBytes)
	if err != nil {
		return err
	}

	w.metricsCollector.Metrics.CallbacksTotal.WithLabelValues("grader callback ").Inc()

	logger.Z(ctx).Debug(ctx, "callback sent", fmt.Sprintf("test for %s", fName), map[string]string{
		"result": string(out),
	})

	return nil
}

func getErrText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
