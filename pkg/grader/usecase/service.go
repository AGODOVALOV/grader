package usecase

import (
	"github.com/AGODOVALOV/grader/pkg/grader/workerpool"
	"github.com/AGODOVALOV/grader/pkg/storage/s3"
	"github.com/AGODOVALOV/grader/pkg/token"
)

type GraderService struct {
	fStorage   *s3.FileStorage
	tokenMaker token.Maker
	WP         *workerpool.WorkerPool
}

func NewGraderService(fStorage *s3.FileStorage, tknMaker token.Maker, wp *workerpool.WorkerPool) *GraderService {
	return &GraderService{
		fStorage:   fStorage,
		tokenMaker: tknMaker,
		WP:         wp,
	}
}
