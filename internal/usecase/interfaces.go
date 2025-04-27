package usecase

import (
	"backend/internal/entity"
	"context"
	"time"
)

type AuthUseCase interface {
	Register(ctx context.Context, email, password, role, name string, DateOfBirth time.Time, profilePhoto, phoneNumber, gender, address, bloodType, rhesus string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GoogleLogin(ctx context.Context, code string) (string, *entity.User, error)
	// VerifyEmail(ctx context.Context, token string) error
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	Logout(ctx context.Context, userID uint, token string) error
}
