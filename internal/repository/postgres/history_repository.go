package postgres

import (
	"backend/internal/entity"
	"context"
	"database/sql"
	"errors"
	"time"
)

type HistoryRepository struct {
	db *sql.DB
}

func NewHistoryRepository(db *sql.DB) *HistoryRepository {
	return &HistoryRepository{
		db: db,
	}
}

func (r *HistoryRepository) Create(ctx context.Context, history *entity.History) error {
	query := `
	INSERT INTO histories (user_id, blood_request_id, next_donation, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)
	`

	now := time.Now()
	history.CreatedAt = now
	history.UpdatedAt = now

	_, err := r.db.ExecContext(
		ctx,
		query,
		history.UserID,
		history.BloodRequestID,
		history.NextDonation,
		history.CreatedAt,
		history.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *HistoryRepository) GetByUserID(ctx context.Context, id uint) ([]*entity.History, error) {
	query := `
	SELECT * FROM histories WHERE user_id = $1;
	`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []*entity.History

	for rows.Next() {
		history := &entity.History{}
		err := rows.Scan(
			&history.ID,
			&history.UserID,
			&history.BloodRequestID,
			&history.ImageDonor,
			&history.NextDonation,
			&history.CreatedAt,
			&history.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}

	return histories, nil
}

func (r *HistoryRepository) NextDonation(ctx context.Context, userID uint) (date time.Time, err error) {
	query := `
	SELECT next_donation FROM histories
	WHERE user_id = $1
	ORDER BY next_donation DESC
	LIMIT 1;
	`

	var history entity.History
	err = r.db.QueryRowContext(ctx, query, userID).Scan(&history.NextDonation)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	return history.NextDonation, nil
}

func (r *HistoryRepository) LatestDonation(ctx context.Context, userID uint) (date time.Time, err error) {
	query := `
	SELECT created_at FROM histories
	WHERE user_id = $1
	ORDER created_at DESC
	LIMIT 1;
	`

	var history entity.History
	err = r.db.QueryRowContext(ctx, query, userID).Scan(&history.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	return history.CreatedAt, nil
}
