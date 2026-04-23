package s3_minio

import (
	"context"
	"mime/multipart"

	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/storage/s3/s3_minio/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(cfg config.Config) (*minio.Client, error) {
	return minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
}

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

func UploadFile(
	ctx context.Context,
	client *minio.Client,
	file multipart.File,
	size int64,
	objectName string,
	bucketName string) error {

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
