package postgres

import (
	"backend/internal/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	INSERT INTO users (email, password, role, name, date_of_birth, profile_photo, phone_number, gender, address, blood_type, rhesus, google_id, fcm_token,created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	RETURNING id
	`

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
		user.Email,        // $1
		user.Password,     // $2
		user.Role,         // $3
		user.Name,         // $4
		user.DateOfBirth,  // $5
		user.ProfilePhoto, // $6
		user.PhoneNumber,  // $7
		user.Gender,       // $8
		user.Address,      // $9
		user.BloodType,    // $10
		user.Rhesus,       // $11
		googleID,          // $12
		user.FCMToken,     // $13
		user.CreatedAt,    // $14
		user.UpdatedAt,    // $15
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

// FindById implements repository.UserRepository.
func (r *UserRepository) FindById(ctx context.Context, id uint) (*entity.User, error) {
	query := `
	SELECT id, email, password, role, name, date_of_birth, profile_photo,
        phone_number, gender, address, blood_type, rhesus, google_id,
		total_donation, coin, fcm_token, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Name,
		&user.DateOfBirth,
		&user.ProfilePhoto,
		&user.PhoneNumber,
		&user.Gender,
		&user.Address,
		&user.BloodType,
		&user.Rhesus,
		&user.GoogleID,
		&user.TotalDonation,
		&user.Coin,
		&user.FCMToken,
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
	query := `
    SELECT 
        id, email, password, role, name, date_of_birth, profile_photo,
        phone_number, gender, address, blood_type, rhesus, google_id,
		total_donation, coin, fcm_token, created_at, updated_at
    FROM users
    WHERE email = $1
    LIMIT 1
    `

	user := &entity.User{}
	var (
		profilePhoto sql.NullString
		googleID     sql.NullString
		fcmToken     sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Name,
		&user.DateOfBirth,
		&profilePhoto,
		&user.PhoneNumber,
		&user.Gender,
		&user.Address,
		&user.BloodType,
		&user.Rhesus,
		&googleID,
		&user.TotalDonation,
		&user.Coin,
		&fcmToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// Set default values untuk field yang tidak di-select
	// user.TotalDonation = 0
	// user.Coin = 0

	// Handle nullable fields
	if profilePhoto.Valid {
		user.ProfilePhoto = &profilePhoto.String
	}
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if fcmToken.Valid {
		user.FCMToken = fcmToken.String
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	return user, nil
}

// FindByGoogleID implements repository.UserRepository.
func (r *UserRepository) FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	query := `
	SELECT id, email, password, role, name, date_of_birth, profile_photo, phone_number, gender, address, blood_type, rhesus, google_id, created_at, updated_at
	FROM users
	WHERE google_id = $1
	`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, googleID).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Name,
		&user.DateOfBirth,
		&user.ProfilePhoto,
		&user.PhoneNumber,
		&user.Gender,
		&user.Address,
		&user.BloodType,
		&user.Rhesus,
		&user.GoogleID,
		&user.TotalDonation,
		&user.Coin,
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
// func (r *UserRepository) VerifyEmail(ctx context.Context, userID uint) error {
// 	query := `
// 	UPDATE users
// 	SET is_email_verified = true, updated_at = $1
// 	WHERE id = $2
// 	`

// 	now := time.Now()

// 	_, err := r.db.ExecContext(
// 		ctx,
// 		query,
// 		now,
// 		userID,
// 	)

// 	return err
// }

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

func (r *UserRepository) UpdateTotalDonation(ctx context.Context, userID uint, totalDonation int) error {
	query := `
		UPDATE users
		SET total_donation = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()

	_, err := r.db.ExecContext(ctx, query, totalDonation, now, userID)
	return err
}

func (r *UserRepository) UpdateCoin(ctx context.Context, userID uint, coin int) error {
	query := `
		UPDATE users
		SET coin = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()

	_, err := r.db.ExecContext(ctx, query, coin, now, userID)
	return err
}

func (r *UserRepository) FindFcmTokenByEmail(ctx context.Context, email string) (string, error) {
	query := `
		SELECT fcm_token
		FROM users
		WHERE email = $1
	`

	var fcmToken string
	err := r.db.QueryRowContext(ctx, query, email).Scan(&fcmToken)
	return fcmToken, err
}

func (r *UserRepository) UpdateFcmTokenByEmail(ctx context.Context, email, fcmToken string) error {
	query := `
		UPDATE users
		SET fcm_token = $1, updated_at = $2
		WHERE email = $3
	`

	now := time.Now()

	_, err := r.db.ExecContext(ctx, query, fcmToken, now, email)
	return err
}

func (r *UserRepository) GetCoinByUserID(ctx context.Context, userID uint) (int, error) {
	var coin int
	query := `
		SELECT coin FROM users WHERE id = $1;
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&coin)
	if err != nil {
		return 0, err
	}
	return coin, nil
}
