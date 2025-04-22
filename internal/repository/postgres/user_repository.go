package postgres

import (
	"backend/internal/entity"
	"context"
	"database/sql"
	"errors"
	"log"
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
	query := `
	INSERT INTO users (email, password, name, google_id, is_email_verified, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	log.Printf("create user repo dipanggil")

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	var googleID interface{}
	if user.GoogleID != nil {
		googleID = *user.GoogleID
	} else {
		googleID = nil
	}

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Password,
		user.Name,
		googleID,
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
	WHERE id = $1
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
			return nil, err
		}
		return nil, err
	}

	return user, nil
}

// FindByEmail implements repository.UserRepository.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
	SELECT id, email, password, name, google_id, is_email_verified, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	user := &entity.User{}
	var googleID sql.NullString

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&googleID,
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

	// Convert NullString ke pointer string
	if googleID.Valid {
		user.GoogleID = &googleID.String
	} else {
		user.GoogleID = nil
	}

	return user, nil
}

// FindByGoogleID implements repository.UserRepository.
func (r *UserRepository) FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	query := `
	SELECT id, email, password, name, google_id, is_email_verified, created_at, updated_at
	FROM users
	WHERE google_id = $1
	`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, googleID).Scan(
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

// Update implements repository.UserRepository.
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
	UPDATE users
	SET email = $1, name = $2, google_id = $3, updated_at = $4
	WHERE id = $5
	`

	user.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Email,
		user.Name,
		user.GoogleID,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

// UpdatePassword implements repository.UserRepository.
func (r *UserRepository) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	query := `
	UPDATE users
	SET password = $1, updated_at = $2
	WHERE id = $3
	`

	now := time.Now()

	_, err := r.db.ExecContext(
		ctx,
		query,
		hashedPassword,
		now,
		userID,
	)

	return err
}

// VerifyEmail implements repository.UserRepository.
func (r *UserRepository) VerifyEmail(ctx context.Context, userID uint) error {
	query := `
	UPDATE users
	SET is_email_verified = true, updated_at = $1
	WHERE id = $2
	`

	now := time.Now()

	_, err := r.db.ExecContext(
		ctx,
		query,
		now,
		userID,
	)

	return err
}

func (r *UserRepository) UpdateProfilePhoto(ctx context.Context, userID uint, photoURL string) error {
	query := `
		UPDATE users
		SET profile_photo = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()

	_, err := r.db.ExecContext(ctx, query, photoURL, now, userID)
	return err
}
