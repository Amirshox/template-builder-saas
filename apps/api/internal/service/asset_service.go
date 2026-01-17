package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"template-builder-api/internal/model"
	"template-builder-api/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type AssetService struct {
	repo        repository.Repository
	minioClient *minio.Client
	bucketName  string
}

func NewAssetService(repo repository.Repository, endpoint, accessKey, secretKey string) (*AssetService, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // For MVP local minio
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init minio: %w", err)
	}

	return &AssetService{
		repo:        repo,
		minioClient: minioClient,
		bucketName:  "assets",
	}, nil
}

func (s *AssetService) EnsureBucket(ctx context.Context) error {
	exists, err := s.minioClient.BucketExists(ctx, s.bucketName)
	if err != nil {
		return err
	}
	if !exists {
		return s.minioClient.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
	}
	return nil
}

func (s *AssetService) UploadAsset(ctx context.Context, orgID uuid.UUID, file io.Reader, filename string, size int64, contentType string) (*model.Asset, error) {
	// 1. Generate unique key
	ext := filepath.Ext(filename)
	assetID := uuid.New()
	s3Key := fmt.Sprintf("%s/%s%s", orgID.String(), assetID.String(), ext)

	// 2. Upload to MinIO
	_, err := s.minioClient.PutObject(ctx, s.bucketName, s3Key, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload to minio: %w", err)
	}

	// 3. Store Metadata
	asset := &model.Asset{
		ID:          assetID,
		OrgID:       orgID,
		Type:        contentType, // Simplify Type mapping
		Filename:    filename,
		ContentType: contentType,
		SizeBytes:   size,
		S3Key:       s3Key,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateAsset(ctx, asset); err != nil {
		return nil, err
	}

	// 4. Generate Presigned URL for response
	url, err := s.minioClient.PresignedGetObject(ctx, s.bucketName, s3Key, 24*time.Hour, nil)
	if err == nil {
		asset.URL = url.String()
	}

	return asset, nil
}

func (s *AssetService) GetDownloadURL(ctx context.Context, assetID uuid.UUID) (string, error) {
	asset, err := s.repo.GetAsset(ctx, assetID)
	if err != nil {
		return "", err
	}

	url, err := s.minioClient.PresignedGetObject(ctx, s.bucketName, asset.S3Key, 1*time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("failed to presign url: %w", err)
	}

	return url.String(), nil
}
