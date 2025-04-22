package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type LocalStorage struct {
	basePath string
	baseURL  string
}

func NewLocalStorage(basePath, baseURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}, nil
}

func (s *LocalStorage) SaveFile(ctx context.Context, fileName string, fileContent io.Reader, contentType string) (*FileInfo, error) {
	// Membuat path file
	filePath := filepath.Join(s.basePath, fileName)

	// membuat direktori jika belum ada
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	// Buat file
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Tulis konten file
	size, err := io.Copy(file, fileContent)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %v", err)
	}

	// Buat URL file
	fileURL := fmt.Sprintf("%s/%s", s.baseURL, fileName)

	return &FileInfo{
		Name:      fileName,
		Size:      size,
		URL:       fileURL,
		CreatedAt: time.Now(),
	}, nil
}

func (s *LocalStorage) GetFile(ctx context.Context, fileName string) (io.ReadCloser, error) {
	filePath := filepath.Join(s.basePath, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	return file, nil
}

func (s *LocalStorage) DeleteFile(ctx context.Context, fileName string) error {
	filePath := filepath.Join(s.basePath, fileName)
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}
