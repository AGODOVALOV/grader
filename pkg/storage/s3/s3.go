package s3

import (
	"context"
	"mime/multipart"

	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/storage/s3/s3_minio"
	"github.com/AGODOVALOV/grader/pkg/storage/s3/s3_minio/config"
	"github.com/minio/minio-go/v7"
)

type FileStorage struct {
	cfg    *config.Config
	Client *minio.Client
}

func NewFileStorage(ctx context.Context, cfg config.Config) (*FileStorage, error) {
	client, err := s3_minio.NewMinioClient(cfg)

	if err != nil {
		logger.Z(ctx).Error(ctx, "NewFileStorage", err.Error())
		return nil, err
	}

	err = s3_minio.EnsureBucket(ctx, client, cfg.Bucket)
	if err != nil {
		logger.Z(ctx).Error(ctx, "check bucket in s3", err.Error())
		return nil, err
	}

	return &FileStorage{
		Client: client,
		cfg:    &cfg,
	}, err
}

func (fs *FileStorage) UploadFile(
	ctx context.Context,
	file multipart.File,
	size int64,
	objectName string) (string, error) {

	err := s3_minio.UploadFile(
		ctx,
		fs.Client,
		file,
		size,
		objectName,
		fs.cfg.Bucket)
	if err != nil {
		return "", err
	}
	return "", nil
}
