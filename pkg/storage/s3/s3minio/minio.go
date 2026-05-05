// Package s3minio provides functionality for interacting with S3.
package s3minio

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/storage/s3/s3minio/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// NewMinioClient creates a new Minio client.
func NewMinioClient(cfg *config.Config) (*minio.Client, error) {
	return minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
}

// EnsureBucket ensures that a bucket exists.
func EnsureBucket(ctx context.Context, client *minio.Client, bucket string) error {
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !exists {
		return client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	}
	return nil
}

// UploadFile uploads a file to a bucket on S3.
func UploadFile(
	ctx context.Context,
	client *minio.Client,
	file multipart.File,
	size int64,
	objectName string,
	bucketName string,
) error {
	info, err := client.PutObject(ctx, bucketName, objectName, file, size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		logger.Z(ctx).Error(ctx, "upload file", err.Error(), map[string]string{
			"bucket":     bucketName,
			"objectName": objectName,
		})
	}

	logger.Z(ctx).Debug(ctx, "upload file", "success", map[string]string{
		"bucket": info.Bucket,
		"key":    info.Key,
	})

	return nil
}

func DownloadFile(ctx context.Context,
	client *minio.Client,
	name string,
	bucket string) ([]byte, error) {

	object, err := client.GetObject(ctx, bucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	defer func(object *minio.Object) {
		err := object.Close()
		if err != nil {
			logger.Z(ctx).Error(ctx, "download file", err.Error(), map[string]string{
				"bucket":   bucket,
				"fileName": name,
			})
		}
	}(object)

	fileBytes, err := io.ReadAll(object)
	if err != nil {
		logger.Z(ctx).Error(ctx, "read file", err.Error(), map[string]string{
			"bucket":   bucket,
			"fileName": name,
		})
		return nil, err
	}

	return fileBytes, nil
}
