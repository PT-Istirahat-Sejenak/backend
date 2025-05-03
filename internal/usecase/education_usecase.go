package usecase

import (
	"backend/internal/entity"
	"backend/internal/infrastructure/storage"
	"backend/internal/repository"
	"context"
)

type educationUseCase struct {
	educationRepo repository.EducationRepository
	fileStorage   storage.FileStorage
}

func NewEducationUseCase(
	educationRepo repository.EducationRepository,
	fileStorage storage.FileStorage,
) EducationUseCase {
	return &educationUseCase{
		educationRepo: educationRepo,
		fileStorage:   fileStorage,
	}
}

// GetAllEducation implements EducationUseCase.
func (e *educationUseCase) GetAllEducation(ctx context.Context) ([]*entity.Education, error) {
	return e.educationRepo.GetAllEducation(ctx)
}

// Post implements EducationUseCase.
func (e *educationUseCase) Post(ctx context.Context, image string, title string, content string, types string) error {
	education := &entity.Education{
		Image:   image,
		Title:   title,
		Content: content,
		Type:    types,
	}

	err := e.educationRepo.Create(ctx, education)

	if err != nil {
		return err
	}

	return nil
}

// FindEducationPencariDonor implements EducationUseCase.
func (e *educationUseCase) FindEducationPencariDonor(ctx context.Context) ([]*entity.Education, error) {
	return e.educationRepo.FindEducationPencariDonor(ctx)
}

// FindEducationPendonor implements EducationUseCase.
func (e *educationUseCase) FindEducationPendonor(ctx context.Context) ([]*entity.Education, error) {
	return e.educationRepo.FindEducationPendonor(ctx)
}

// FindById implements EducationUseCase.
func (e *educationUseCase) FindById(ctx context.Context, id uint) (*entity.Education, error) {
	return e.educationRepo.FindById(ctx, id)
}

// Update implements EducationUseCase.
func (e *educationUseCase) Update(ctx context.Context, image string, title string, content string, types string) error {
	education := &entity.Education{
		Image:   image,
		Title:   title,
		Content: content,
		Type:    types,
	}

	err := e.educationRepo.Update(ctx, education)

	if err != nil {
		return err
	}

	return nil
}

// Delete implements EducationUseCase.
func (e *educationUseCase) Delete(ctx context.Context, id uint) error {
	return e.educationRepo.Delete(ctx, id)
}
