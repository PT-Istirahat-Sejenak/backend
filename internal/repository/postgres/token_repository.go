package postgres

import (
	"backend/internal/entity"
	"context"
	"database/sql"
	"errors"
)

type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

func (r *TokenRepository) Create(ctx context.Context, token *entity.Token) error {
	query := `
	INSERT INTO tokens (user_id, token, type, expires_at, created_at) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		token.UserID,
		token.Token,
		token.Type,
		token.ExpiresAt,
		token.CreatedAt,
	).Scan(&token.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *TokenRepository) FindByToken(ctx context.Context, token string, tokenType entity.TokenType) (*entity.Token, error) {
	query := `
	SELECT id, user_id, token, type, expires_at, created_at
	FROM tokens
	WHERE token = $1 AND type = $2
	`

	t := &entity.Token{}
	err := r.db.QueryRowContext(ctx, query, token, tokenType).Scan(
		&t.ID,
		&t.UserID,
		&t.Token,
		&t.Type,
		&t.ExpiresAt,
		&t.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return t, nil
}

func (r *TokenRepository) Delete(ctx context.Context, id uint) error {
	query := `
	DELETE FROM tokens WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

func (r *TokenRepository) FindByUserID(ctx context.Context, userID uint) ([]*entity.Token, error) {
	query := `
        SELECT id, user_id, token, type, expires_at, created_at 
        FROM tokens 
        WHERE user_id = $1
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*entity.Token
	for rows.Next() {
		var t entity.Token
		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Token,
			&t.Type,
			&t.ExpiresAt,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &t)
	}

	return tokens, nil
}

func (r *TokenRepository) Exists(ctx context.Context, token string, tokenType entity.TokenType) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM tokens 
            WHERE token = $1 AND type = $2 AND expires_at > NOW()
        )
    `

	var exists bool
	err := r.db.QueryRowContext(ctx, query, token, tokenType).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
