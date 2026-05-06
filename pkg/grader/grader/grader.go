package grader

import (
	"context"

	"github.com/AGODOVALOV/grader/pkg/config/config"
	"github.com/AGODOVALOV/grader/pkg/grader/handler"
	"github.com/AGODOVALOV/grader/pkg/grader/workerpool"
	"github.com/AGODOVALOV/grader/pkg/storage/s3"
	"github.com/AGODOVALOV/grader/pkg/token"
)

type Grader struct {
	Handler *handler.GraderHandler
}

func NewGrader(ctx context.Context, fStorage *s3.FileStorage, tknMaker token.Maker, cfg *config.Config) *Grader {
	return &Grader{
		Handler: handler.NewGraderHandler(fStorage, tknMaker, workerpool.NewWorkerPool(ctx, cfg, fStorage)),
	}
}
