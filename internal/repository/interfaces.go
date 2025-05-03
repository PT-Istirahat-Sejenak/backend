package repository

import (
	"backend/internal/entity"
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindById(ctx context.Context, id uint) (*entity.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error
	UpdateProfilePhoto(ctx context.Context, userID uint, photoURL string) error
	// VerifyEmail(ctx context.Context, userID uint) error
}

type TokenRepository interface {
	Create(ctx context.Context, token *entity.Token) error
	FindByToken(ctx context.Context, token string, tokenType entity.TokenType) (*entity.Token, error)
	Delete(ctx context.Context, id uint) error
	FindByUserID(ctx context.Context, userID uint) ([]*entity.Token, error)
	Exists(ctx context.Context, token string, tokenType entity.TokenType) (bool, error)
}

type EducationRepository interface {
	Create(ctx context.Context, education *entity.Education) error
	GetAllEducation(ctx context.Context) ([]*entity.Education, error)
	FindEducationPendonor(ctx context.Context) ([]*entity.Education, error)
	FindEducationPencariDonor(ctx context.Context) ([]*entity.Education, error)
	FindById(ctx context.Context, id uint) (*entity.Education, error)
	Update(ctx context.Context, education *entity.Education) error
	Delete(ctx context.Context, id uint) error
}

type UploadEvidence interface {
	Upload(ctx context.Context, uploadEvidence *entity.UploadEvidence) error
}

type Histories interface {
	Create(ctx context.Context, history *entity.History) error
	GetByUserID(ctx context.Context, userID uint) ([]*entity.History, error)
	NextDonation(ctx context.Context, userID uint) (date time.Time, err error)
	LatestDonation(ctx context.Context, userID uint) (date time.Time, err error)
}
