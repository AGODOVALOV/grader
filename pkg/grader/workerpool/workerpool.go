package workerpool

import (
	"context"
	"fmt"
	"sync"

	"github.com/AGODOVALOV/grader/pkg/config/config"
	"github.com/AGODOVALOV/grader/pkg/dto"
	graderconfig "github.com/AGODOVALOV/grader/pkg/grader/config"
	"github.com/AGODOVALOV/grader/pkg/logger"
)

type WorkerPool struct {
	cfg   *graderconfig.Config
	Tasks chan *dto.GraderPayload
	wg    *sync.WaitGroup
}

func NewWorkerPool(cfg *config.Config) *WorkerPool {
	return &WorkerPool{
		Tasks: make(chan *dto.GraderPayload, cfg.Grader.Workers*20),
		cfg:   &cfg.Grader,
		wg:    &sync.WaitGroup{},
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

			logger.Z(ctx).Debug(ctx, "process new task", fmt.Sprintf("%v", v))
		}
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
