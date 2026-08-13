package minio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"match-me-api/internal/config"
	"match-me-api/internal/logger"
	"net/url"
	"path"
	"strings"

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
		Creds:  credentials.NewStaticV4(minioCfg.RootUser, minioCfg.RootPassword, ""),
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

// InitBucket creates bucket in MinIO if it doesn't exist and applies the policy.
func (s *Storage) InitBucket(ctx context.Context, bucketName, policyJSON string) error {
	const op = "storage.minio.InitBucket"

	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("%s: failed to check if bucket exists '%s': %w", op, bucketName, err)
	}

	if !exists {
		if err := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("%s: failed to create bucket '%s': %w", op, bucketName, err)
		}
		s.log.Info("s3 bucket created successfully", slog.String("bucket", bucketName))
	} else {
		s.log.Info("s3 bucket already exists", slog.String("bucket", bucketName))
	}

	if policyJSON != "" {
		if err := s.client.SetBucketPolicy(ctx, bucketName, policyJSON); err != nil {
			return fmt.Errorf("%s: failed to set policy for bucket '%s': %w", op, bucketName, err)
		}
		s.log.Info("s3 bucket policy applied successfully", slog.String("bucket", bucketName))
	}

	return nil
}

// ObjectExists checks whether an object exists in the specified MinIO bucket.
func (s *Storage) ObjectExists(ctx context.Context, bucketName, objectName string) (bool, error) {
	const op = "storage.minio.ObjectExists"

	_, err := s.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("%s: failed to stat object '%s' in bucket '%s': %w", op, objectName, bucketName, err)
	}
	return true, nil
}

// PublicURL constructs the public URL for an object in a bucket.
func (s *Storage) PublicURL(bucketName, objectName string) string {
	const op = "storage.minio.PublicURL"

	if s.baseURL == "" {
		return objectName
	}

	u, err := url.Parse(s.baseURL)
	if err != nil {
		s.log.Error("failed to parse base url", slog.String("op", op), slog.String("url", s.baseURL), logger.Err(err))
		return objectName
	}

	u.Path = path.Join(u.Path, bucketName, objectName)
	return u.String()
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

	s.log.Info("S3 object uploaded successfully", slog.String("bucket", bucketName), slog.String("object", objectName))
	return u.String(), nil
}

// Delete removes an object from specified MinIO bucket.
func (s *Storage) Delete(ctx context.Context, bucketName, objectName string) error {
	const op = "storage.minio.Delete"

	if objectName == "" {
		return nil
	}

	err := s.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("%s: failed to delete object '%s': %w", op, objectName, err)
	}

	s.log.Info("S3 object deleted successfully", slog.String("bucket", bucketName), slog.String("object", objectName))
	return nil
}

// MakeReadOnlyPolicy generates a standard S3 public read-only policy JSON for a bucket.
func MakeReadOnlyPolicy(bucketName string) (string, error) {
	const op = "storage.minio.MakeReadOnlyPolicy"

	policy := map[string]any{
		"Version": "2012-10-17",
		"Statement": []map[string]any{
			{
				"Effect":    "Allow",
				"Principal": "*",
				"Action":    []string{"s3:GetObject"},
				"Resource":  []string{fmt.Sprintf("arn:aws:s3:::%s/*", bucketName)},
			},
		},
	}

	data, err := json.Marshal(policy)
	if err != nil {
		return "", fmt.Errorf("%s: failed to marshal policy: %w", op, err)
	}

	return string(data), nil
}

// ExtractObjectName extracts object name from public URL saved in DB
func (s *Storage) ExtractObjectName(rawURL, bucketName string) string {
	const op = "storage.minio.ExtractObjectName"

	if rawURL == "" {
		return ""
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	// URL: http://localhost:9000/avatars/users/42/1710000000.png
	// u.Path: /avatars/users/42/1710000000.png
	cleanPath := strings.TrimPrefix(u.Path, "/")
	prefix := bucketName + "/"

	if strings.HasPrefix(cleanPath, prefix) {
		return strings.TrimPrefix(cleanPath, prefix)
	}

	return ""
}
