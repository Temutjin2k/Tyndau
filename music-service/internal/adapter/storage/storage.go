package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioConfig defines required configuration for MinIO connection.
type MinioConfig struct {
	ADDR      string `env:"MINIO_ADDR,required"`       // e.g. "localhost:9000"
	AccessKey string `env:"MINIO_ACCESS_KEY,required"` // Access key for MinIO
	SecretKey string `env:"MINIO_SECRET_KEY,required"` // Secret key for MinIO
	Bucket    string `env:"MINIO_BUCKET,required"`     // e.g. "song"
	UseSSL    bool   `env:"MINIO_USE_SSL"`             // true if HTTPS, false otherwise
	PublicURL string `env:"MINIO_PUBLIC_URL,required"` // e.g. "https://cdn.yourdomain.com"
}

// MinioStorage is an implementation of the Uploader interface using MinIO.
type MinioStorage struct {
	client    *minio.Client
	bucket    string
	publicURL string
}

// NewMinioStorage initializes MinIO client with given configuration.
func NewMinioStorage(ctx context.Context, cfg MinioConfig) (*MinioStorage, error) {
	client, err := minio.New(cfg.ADDR, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	// Ensure bucket exists or create it
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &MinioStorage{
		client:    client,
		bucket:    cfg.Bucket,
		publicURL: cfg.PublicURL,
	}, nil
}

// PresignedPutURL returns a presigned URL to upload an object to MinIO.
func (m *MinioStorage) PresignedPutURL(ctx context.Context, bucket, objectName string, expires time.Duration) (string, error) {
	url, err := m.client.PresignedPutObject(ctx, bucket, objectName, expires)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned PUT URL: %w", err)
	}

	return url.String(), nil
}

// GenerateSongURLs returns both the presigned PUT upload URL and the public GET URL
func (m *MinioStorage) GenerateSongURLs(ctx context.Context, songID int64, fileName string) (string, string, error) {
	// Create object key, e.g. "songs/1234_filename.mp3"
	objectKey := fmt.Sprintf("songs/%d_%s", songID, fileName)

	// Generate presigned PUT URL
	uploadURL, err := m.client.PresignedPutObject(ctx, m.bucket, objectKey, time.Minute*10)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	// Build public URL manually using public base URL
	publicURL := fmt.Sprintf("%s/%s/%s", m.publicURL, m.bucket, objectKey)

	return uploadURL.String(), publicURL, nil
}
