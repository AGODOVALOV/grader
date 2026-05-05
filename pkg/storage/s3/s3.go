// Package s3 provides a file storage service using S3
package s3

import (
	"context"
	"mime/multipart"

	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/storage/s3/s3minio"
	"github.com/AGODOVALOV/grader/pkg/storage/s3/s3minio/config"
	"github.com/minio/minio-go/v7"
)

// FileStorage represents a file storage service.
type FileStorage struct {
	cfg    *config.Config
	Client *minio.Client
}

// NewFileStorage creates a new instance of FileStorage.
func NewFileStorage(ctx context.Context, cfg *config.Config) (*FileStorage, error) {
	client, err := s3minio.NewMinioClient(cfg)
	if err != nil {
		logger.Z(ctx).Error(ctx, "NewFileStorage", err.Error())
		return nil, err
	}

	err = s3minio.EnsureBucket(ctx, client, cfg.Bucket)
	if err != nil {
		logger.Z(ctx).Error(ctx, "check bucket in s3", err.Error())
		return nil, err
	}

	return &FileStorage{
		Client: client,
		cfg:    cfg,
	}, nil
}

// UploadFile uploads a file to S3.
func (fs *FileStorage) UploadFile(
	ctx context.Context,
	file multipart.File,
	size int64,
	objectName string,
) (string, error) {
	err := s3minio.UploadFile(
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

func (fs *FileStorage) DownloadFile(
	ctx context.Context,
	objectName string) ([]byte, error) {
	fileData, err := s3minio.DownloadFile(
		ctx,
		fs.Client,
		objectName,
		fs.cfg.Bucket)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}
