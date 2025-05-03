package usecase

import (
	"backend/internal/entity"
	"backend/internal/infrastructure/storage"
	"backend/internal/repository"
	"context"
)

type uploadEvidenceUseCase struct {
	uploadEvidenceRepo repository.UploadEvidence
	fileStorage        storage.FileStorage
}

func NewUploadEvidenceUseCase(
	uploadEvidenceRepo repository.UploadEvidence,
	fileStorage storage.FileStorage,
) EvidenceUseCase {
	return &uploadEvidenceUseCase{
		uploadEvidenceRepo: uploadEvidenceRepo,
		fileStorage:        fileStorage,
	}
}

// UploadEvidence implements EvidenceUseCase.
func (u *uploadEvidenceUseCase) UploadEvidence(ctx context.Context, userID uint, photoURL string) error {
	upload := &entity.UploadEvidence{
		UserID: userID,
		Image:  photoURL,
	}

	err := u.uploadEvidenceRepo.Upload(ctx, upload)
	return err
}
