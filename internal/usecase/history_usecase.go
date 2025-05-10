package usecase

import (
	"backend/internal/entity"
	"backend/internal/infrastructure/storage"
	"backend/internal/repository"
	"context"
	"time"
)

type historyUseCase struct {
	historyRepo repository.HistoriesRepository
	fileStorage storage.FileStorage
}

func NewHistoryUseCase(historyRepo repository.HistoriesRepository, fileStorage storage.FileStorage) HistoryUseCase {
	return &historyUseCase{
		historyRepo: historyRepo,
		fileStorage: fileStorage,
	}
}

// AddHistory implements HistoryUseCase.
func (h *historyUseCase) AddHistory(ctx context.Context, userID uint, bloodRequestID uint, imageDonor string, nextDonation time.Time) error {
	history := &entity.History{
		UserID:         userID,
		BloodRequestID: bloodRequestID,
		ImageDonor:     imageDonor,
		NextDonation:   nextDonation,
	}

	err := h.historyRepo.Create(ctx, history)

	if err != nil {
		return err
	}

	return nil
}

// HistoryByUserId implements HistoryUseCase.
func (h *historyUseCase) HistoryByUserId(ctx context.Context, userID uint) ([]*entity.History, error) {
	return h.historyRepo.GetByUserID(ctx, userID)
}

// GetLatestDonation implements HistoryUseCase.
func (h *historyUseCase) GetLatestDonation(ctx context.Context, userID uint) (date time.Time, err error) {
	return h.historyRepo.LatestDonation(ctx, userID)
}

// GetNextDonation implements HistoryUseCase.
func (h *historyUseCase) GetNextDonation(ctx context.Context, userID uint) (date time.Time, err error) {
	return h.historyRepo.NextDonation(ctx, userID)
}
