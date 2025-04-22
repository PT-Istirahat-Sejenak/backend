package storage

import (
	"context"
	"io"
	"time"
)

// FileInfo berisi informasi tentang file yang disimpan
type FileInfo struct {
	Name      string
	Size      int64
	URL       string
	CreatedAt time.Time
}

// FileStorage mendefinisikan interface untuk operasi penyimpanan file
type FileStorage interface {
	// SaveFile menyimpan file dan mengembalikan informasi file
	SaveFile(ctx context.Context, fileName string, fileContent io.Reader, contentType string) (*FileInfo, error)
	
	// GetFile mengembalikan file sebagai io.ReadCloser
	GetFile(ctx context.Context, fileName string) (io.ReadCloser, error)
	
	// DeleteFile menghapus file
	DeleteFile(ctx context.Context, fileName string) error
}