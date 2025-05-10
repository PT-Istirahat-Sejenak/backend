package usecase

import (
	"backend/internal/entity"
	"backend/internal/infrastructure/oauth"
	"context"
	"mime/multipart"
	"time"
)

type AuthUseCase interface {
	Register(ctx context.Context, email, password, role, name string, DateOfBirth time.Time, profilePhoto, phoneNumber, gender, address, bloodType, rhesus, fcmToken string) (*entity.User, string, error)
	Login(ctx context.Context, email, password string) (string, error)
	GoogleLogin(ctx context.Context, code string) (string, *entity.User, error)
	GoogleLoginMobile(ctx context.Context, googleInfo oauth.GoogleUserInfo) (string, *entity.User, error)
	ValidateFcmToken(ctx context.Context, userEmail, fcmToken string) error
	// VerifyEmail(ctx context.Context, token string) error
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	UpdateCountDonation(ctx context.Context, userID uint, total int) error
	UpdateCoinTotal(ctx context.Context, userID uint, coin int) error
	Logout(ctx context.Context, userID uint, token string) error
}

type ProfileUseCase interface {
	UpdateProfilePhoto(ctx context.Context, userID uint, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	GetUserProfile(ctx context.Context, userID uint) (*entity.User, error)
}

type EducationUseCase interface {
	Post(ctx context.Context, image, title, content, types string) error
	GetAllEducation(ctx context.Context) ([]*entity.Education, error)
	FindEducationPendonor(ctx context.Context) ([]*entity.Education, error)
	FindEducationPencariDonor(ctx context.Context) ([]*entity.Education, error)
	FindById(ctx context.Context, id uint) (*entity.Education, error)
	Update(ctx context.Context, image, title, content, types string) error
	Delete(ctx context.Context, id uint) error
}

type EvidenceUseCase interface {
	UploadEvidence(ctx context.Context, userID uint, photoURL string) error
}

type HistoryUseCase interface {
	AddHistory(ctx context.Context, userID, bloodRequestID uint, imageDonor string, nextDonation time.Time) error
	HistoryByUserId(ctx context.Context, userID uint) ([]*entity.History, error)
	GetNextDonation(ctx context.Context, userID uint) (date time.Time, err error)
	GetLatestDonation(ctx context.Context, userID uint) (date time.Time, err error)
}

type ChatbotUsecase interface {
	GetReply(ctx context.Context, userID uint, message string) (string, error)
}

type RewardUseCase interface {
	GetReward(ctx context.Context, userID uint, operatorID, amount, number string) (res string, err error)
	GetBalance(ctx context.Context) (res float64, err error)
	GetToken(ctx context.Context) (res *entity.TokenReward, err error)
}

type FCMUseCase interface {
	GetAccessToken(ctx context.Context) (*entity.Fcm, error)
	SendFCMV1(ctx context.Context, userID uint, title, body string) error
}
