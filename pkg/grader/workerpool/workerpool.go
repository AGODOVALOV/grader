package workerpool

import (
	"context"
	"sync"

	"github.com/AGODOVALOV/grader/pkg/config/config"
	"github.com/AGODOVALOV/grader/pkg/dto"
	graderconfig "github.com/AGODOVALOV/grader/pkg/grader/config"
	"github.com/AGODOVALOV/grader/pkg/grader/workerpool/worker"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/storage/s3"
)

type WorkerPool struct {
	cfg    *graderconfig.Config
	Tasks  chan *dto.GraderPayload
	wg     *sync.WaitGroup
	worker *worker.Worker
}

func NewWorkerPool(cfg *config.Config, fStorage *s3.FileStorage) *WorkerPool {
	return &WorkerPool{
		Tasks:  make(chan *dto.GraderPayload, cfg.Grader.Workers*20),
		cfg:    &cfg.Grader,
		wg:     &sync.WaitGroup{},
		worker: worker.NewWorker(fStorage),
	}
}

func (wp *WorkerPool) StartProcessingGradeTasks(ctx context.Context) {
	for range wp.cfg.Workers {
		wp.wg.Go(func() {
			wp.ProcessTask(ctx, wp.Tasks)
		})
	}
}

func (wp *WorkerPool) ProcessTask(ctx context.Context, ch <-chan *dto.GraderPayload) {
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-ch:
			if !ok {
				return
			}
			err := wp.worker.DoJob(ctx, v)
			if err != nil {
				logger.Z(ctx).Error(ctx, "grader processing new task", err.Error(), map[string]string{
					"userID":   v.UserID,
					"reviewID": v.ReviewID,
				})
			}
		}
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
