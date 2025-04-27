package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client     *s3.Client
	bucketName string
	baseURL    string
}

func NewS3Storage(accountID, accessKeyID, accessKeySecret, bucketName, region, baseURL string) (*S3Storage, error) {
	endpointURL := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	// Buat custom resolver untuk endpoint
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: endpointURL,
		}, nil
	})

	// Konfigurasi kredensial
	creds := credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")

	// Buat konfigurasi AWS
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	if err != nil {
		return nil, fmt.Errorf("error loading AWS configuration: %w", err)
	}

	// Buat klien S3
	client := s3.NewFromConfig(cfg)

	return &S3Storage{
		client:     client,
		bucketName: bucketName,
		baseURL:    baseURL,
	}, nil
}

func (s *S3Storage) SaveFile(ctx context.Context, fileName string, fileContent io.Reader, contentType string) (*FileInfo, error) {
	// Buat request untuk upload file
	putObjectInput := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(fileName),
		Body:        fileContent,
		ContentType: aws.String(contentType),
	}

	// Upload file
	_, err := s.client.PutObject(ctx, putObjectInput)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}

	// Buat URL file
	fileURL := fmt.Sprintf("%s/%s", s.baseURL, fileName)

	return &FileInfo{
		Name:      fileName,
		Size:      0, // Tidak ada informasi ukuran file pada s3
		URL:       fileURL,
		CreatedAt: time.Now(),
	}, nil
}

func (s *S3Storage) GetFile(ctx context.Context, fileName string) (io.ReadCloser, error) {
	// Buat request untuk download file
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	}

	// Download file
	output, err := s.client.GetObject(ctx, getObjectInput)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}

	return output.Body, nil
}

func (s *S3Storage) DeleteFile(ctx context.Context, fileName string) error {
	// Buat request untuk hapus file
	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	}

	// Hapus file
	_, err := s.client.DeleteObject(ctx, deleteObjectInput)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}
