package minio

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"match-me-api/internal/config"
	"net/url"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	client  *minio.Client
	baseURL string
	log     *slog.Logger
}

func NewMinioStorage(ctx context.Context, minioCfg config.MinIOConfig, log *slog.Logger) (*Storage, error) {
	const op = "storage.minio.New"

	client, err := minio.New(minioCfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioCfg.AccessKeyID, minioCfg.SecretKey, ""),
		Secure: minioCfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to init minio client: %w", op, err)
	}

	return &Storage{
		client:  client,
		baseURL: minioCfg.PublicEndpoint,
		log:     log,
	}, nil
}

// InitBucket creates bucket for minio
func (s *Storage) InitBucket(ctx context.Context, bucketName string) error {
	const op = "storage.minio.InitBucket"

	bucket, err := s.client.BucketExists(ctx, bucketName)
	if err == nil && bucket {
		return nil
	}

	if err := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
		bucketExists, checkErr := s.client.BucketExists(ctx, bucketName)
		if checkErr == nil && bucketExists {
			return nil
		}
		return fmt.Errorf("%s: failed to create bucket '%s': %w", op, bucketName, err)
	}

	return nil
}

// Upload saves data stream into MinIO with Content-Type identifier and return URL to it.
func (s *Storage) Upload(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	const op = "storage.minio.Upload"

	_, err := s.client.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("%s: failed to upload object '%s': %w", op, objectName, err)
	}

	if s.baseURL == "" {
		return objectName, nil
	}

	u, err := url.Parse(s.baseURL)
	if err != nil {
		return "", fmt.Errorf("%s: invalid base url: %w", op, err)
	}

	u.Path = path.Join(u.Path, bucketName, objectName)
	return u.String(), nil
}
