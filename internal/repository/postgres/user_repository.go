package postgres

import (
	"backend/internal/entity"
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (email, password, name, google_id, is_email_verified, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	Returning id
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Password,
		user.Name,
		user.GoogleID,
		user.IsEmailVerified,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

// FindById implements repository.UserRepository.
func (r *UserRepository) FindById(ctx context.Context, id uint) (*entity.User, error) {
	query := `
	SELECT id, email, password, name, google_id, is_email_verified, created_at, updated_at	
	FROM users
	WHERE email = $1
	`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.GoogleID,
		&user.IsEmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// FindByEmail implements repository.UserRepository.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	panic("unimplemented")
}

// FindByGoogleID implements repository.UserRepository.
func (r *UserRepository) FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	panic("unimplemented")
}

// Update implements repository.UserRepository.
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	panic("unimplemented")
}

// UpdatePassword implements repository.UserRepository.
func (r *UserRepository) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	panic("unimplemented")
}

// VerifyEmail implements repository.UserRepository.
func (r *UserRepository) VerifyEmail(ctx context.Context, userID uint) error {
	panic("unimplemented")
}
