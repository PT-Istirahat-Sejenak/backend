package postgres

import (
	"context"
	"database/sql"
)

type RewardRepository struct {
	db *sql.DB
}

func NewRewardRepository(db *sql.DB) *RewardRepository {
	return &RewardRepository{
		db: db,
	}
}

func (r *RewardRepository) GetPriceByAmount(ctx context.Context, amount int) (int, error) {
	var price int
	query := `
		SELECT price FROM rewards WHERE amount = $1;
	`
	err := r.db.QueryRowContext(ctx, query, amount).Scan(&price)
	return price, err
}
