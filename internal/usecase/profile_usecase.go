package usecase

import (
	"backend/internal/entity"
	"backend/internal/infrastructure/storage"
	"backend/internal/repository"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type profileUseCase struct {
	userRepo    repository.UserRepository
	fileStorage storage.FileStorage
}

func NewProfileUseCase(
	userRepo repository.UserRepository, fileStorage storage.FileStorage,
) ProfileUseCase {
	return &profileUseCase{
		userRepo:    userRepo,
		fileStorage: fileStorage,
	}
}

func (p *profileUseCase) UpdateProfilePhoto(
	ctx context.Context,
	userID uint,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) (string, error) {
	// Validasi file
	if fileHeader.Size > 5*1024*1024 { // 5MB
		return "", fmt.Errorf("file size exceeds maximum limit of 5MB")
	}

	// Validasi tipe file (hanya terima gambar)
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	if !allowedExts[ext] {
		return "", fmt.Errorf("invalid file type, only JPG, JPEG, PNG, and GIF are allowed")
	}

	// Buat nama file yang unik
	timestamp := time.Now().UnixNano()
	fileName := fmt.Sprintf("profiles/%d/%d%s", userID, timestamp, ext)

	// Tentukan content type berdasarkan ekstensi file
	contentType := ""
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	}

	// Simpan file
	fileInfo, err := p.fileStorage.SaveFile(ctx, fileName, file, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to save profile photo: %v", err)
	}

	// Update profil user di database
	err = p.userRepo.UpdateProfilePhoto(ctx, userID, fileInfo.URL)
	if err != nil {
		// Jika update database gagal, hapus file yang sudah diupload
		_ = p.fileStorage.DeleteFile(ctx, fileName)
		return "", fmt.Errorf("failed to update user profile: %v", err)
	}

	return fileInfo.URL, nil
}

func (p *profileUseCase) GetUserProfile(ctx context.Context, userID uint) (*entity.User, error) {
	return p.userRepo.FindById(ctx, userID)
}
