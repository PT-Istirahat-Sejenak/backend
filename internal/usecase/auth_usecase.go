package usecase

import (
	"backend/internal/entity"
	"backend/internal/infrastructure/email"
	"backend/internal/infrastructure/oauth"
	"backend/internal/repository"
	"backend/pkg/hash"
	"backend/pkg/jwt"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"
)

type authUseCase struct {
	userRepo     repository.UserRepository
	tokenRepo    repository.TokenRepository
	jwtService   *jwt.JWTService
	emailService *email.EmailService
	googleOauth  *oauth.GooogleOauth
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtService *jwt.JWTService,
	emailService *email.EmailService,
	googleOauth *oauth.GooogleOauth,
) AuthUseCase {
	return &authUseCase{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		jwtService:   jwtService,
		emailService: emailService,
		googleOauth:  googleOauth,
	}
}

func (a *authUseCase) generateOTP(length int) (string, error) {
	const digits = "0123456789"
	b := make([]byte, length)
	for i := range b {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		b[i] = digits[randomIndex.Int64()]
	}

	return string(b), nil
}

// Register implements AuhtUseCase.
func (a *authUseCase) Register(ctx context.Context, email, password, role, name string, DateOfBirth time.Time, profilePhoto, phoneNumber, gender, address, bloodType, rhesus, fcmToken string) (*entity.User, string, error) {
	existingUser, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "",err
	}

	if existingUser != nil {
		return nil, "",errors.New("user with this email already exists")
	}

	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		return nil, "",err
	}

	user := &entity.User{
		Email:        email,
		Password:     hashedPassword,
		Name:         name,
		Role:         role,
		DateOfBirth:  DateOfBirth,
		ProfilePhoto: &profilePhoto,
		PhoneNumber:  phoneNumber,
		Gender:       gender,
		Address:      address,
		BloodType:    bloodType,
		Rhesus:       rhesus,
		GoogleID:     nil,
	}

	err = a.userRepo.Create(ctx, user)
	if err != nil {
		return nil, "", err
	}

	token, err := a.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login implements AuhtUseCase.
func (a *authUseCase) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("invalid email or password")
	}

	if !hash.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := a.jwtService.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GoogleLogin implements AuhtUseCase.
func (a *authUseCase) GoogleLogin(ctx context.Context, code string) (string, *entity.User, error) {
	googleUser, err := a.googleOauth.GetUserInfo(ctx, code)
	if err != nil {
		return "", nil, err
	}

	// Check if user already exists by google id
	user, err := a.userRepo.FindByGoogleID(ctx, googleUser.ID)
	if err != nil {
		return "", nil, err
	}

	var isNewUser bool

	if user == nil {
		// Check if user exist by email
		user, err = a.userRepo.FindByEmail(ctx, googleUser.Email)
		if err != nil {
			return "", nil, err
		}

		if user == nil {
			// // Create new user
			// user = entity.NewGoogleUser(
			// 	googleUser.Email,
			// 	googleUser.Name,
			// 	googleUser.ID,
			// 	true,
			// )
			googleID := googleUser.ID

			user = &entity.User{
				Email:        googleUser.Email,
				Name:         googleUser.Name,
				ProfilePhoto: &googleUser.Picture,
				GoogleID:     &googleID,
			}

			err = a.userRepo.Create(ctx, user)
			if err != nil {
				return "", nil, err
			}

			isNewUser = true
		} else {
			// Update existing user with Google ID

			googleID := googleUser.ID
			user.GoogleID = &googleID

			err = a.userRepo.Update(ctx, user)
			if err != nil {
				return "", nil, err
			}
		}
	}

	// Generate JWT token
	token, err := a.jwtService.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	if isNewUser {
		return token, user, nil
	}

	return token, nil, nil
}

// GoogleLoginMobile implements AuthUseCase.
func (a *authUseCase) GoogleLoginMobile(ctx context.Context, googleInfo oauth.GoogleUserInfo) (string, *entity.User, error) {
	user, err := a.userRepo.FindByGoogleID(ctx, googleInfo.ID)
	if err != nil {
		return "", nil, err
	}

	var isNewUser bool

	if user == nil {
		user, err = a.userRepo.FindByEmail(ctx, googleInfo.Email)
		if err != nil {
			return "", nil, err
		}

		if user == nil {
			googleID := googleInfo.ID

			user = &entity.User{
				Email:        googleInfo.Email,
				Name:         googleInfo.Name,
				ProfilePhoto: &googleInfo.Picture,
				GoogleID:     &googleID,
			}

			err = a.userRepo.Create(ctx, user)
			if err != nil {
				return "", nil, err
			}

			isNewUser = true
		} else {
			googleID := googleInfo.ID
			user.GoogleID = &googleID

			err = a.userRepo.Update(ctx, user)
			if err != nil {
				return "", nil, err
			}
		}
	}

	// Generate JWT token
	token, err := a.jwtService.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	if isNewUser {
		return token, user, nil
	}

	return token, nil, nil
}

// // VerifyEmail implements AuhtUseCase.
// func (a *authUseCase) VerifyEmail(ctx context.Context, token string) error {
// 	tokenEntity, err := a.tokenRepo.FindByToken(ctx, token, entity.EmailVerify)
// 	if err != nil {
// 		return err
// 	}

// 	if tokenEntity == nil {
// 		return errors.New("invalid or expired token")
// 	}

// 	if tokenEntity.ExpiresAt.Before(time.Now()) {
// 		return errors.New("token expired")
// 	}

// 	err = a.userRepo.VerifyEmail(ctx, tokenEntity.UserID)
// 	if err != nil {
// 		return err
// 	}

// 	return a.tokenRepo.Delete(ctx, tokenEntity.ID)
// }

// RequestPasswordReset implements AuhtUseCase.
func (a *authUseCase) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user with this email does not exist")
	}

	// Generate reset password token
	resetCode, err := a.generateOTP(6)
	if err != nil {
		return err
	}

	token := &entity.Token{
		UserID:    user.ID,
		Token:     resetCode,
		Type:      entity.ResetPassword,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	err = a.tokenRepo.Create(ctx, token)
	if err != nil {
		return err
	}

	// Send reset password email
	// return a.emailService.SendResetPasswordEmail(user.Email, resetCode)
	err = a.emailService.SendResetPasswordEmail(user.Email, resetCode)
	if err != nil {
		log.Printf("Failed to send verification email :%v", err)
		return fmt.Errorf("failed to send verification email")
	}

	return nil
}

// ResetPassword implements AuhtUseCase.
func (a *authUseCase) ResetPassword(ctx context.Context, token string, newPassword string) error {
	tokenEntity, err := a.tokenRepo.FindByToken(ctx, token, entity.ResetPassword)
	if err != nil {
		return err
	}

	if tokenEntity == nil {
		return errors.New("invalid or expired token")
	}

	if tokenEntity.ExpiresAt.Before(time.Now()) {
		return errors.New("token has expired")
	}

	hashedPassword, err := hash.HashPassword(newPassword)
	if err != nil {
		return err
	}

	err = a.userRepo.UpdatePassword(ctx, tokenEntity.UserID, hashedPassword)
	if err != nil {
		return err
	}

	return a.tokenRepo.Delete(ctx, tokenEntity.ID)
}

// GetUserByID implements AuhtUseCase.
func (a *authUseCase) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	return a.userRepo.FindById(ctx, id)
}

func (a *authUseCase) UpdateCountDonation(ctx context.Context, userID uint, totalDonation int) error {
	return a.userRepo.UpdateTotalDonation(ctx, userID, totalDonation)
}

func (a *authUseCase) UpdateCoinTotal(ctx context.Context, userID uint, coin int) error {
	return a.userRepo.UpdateCoin(ctx, userID, coin)
}

func (a *authUseCase) Logout(ctx context.Context, userID uint, token string) error {
	// Validasi token
	claims, err := a.jwtService.ValidateToken(token)
	if err != nil {
		return err
	}

	// Pastikan token milik user yang sesuai
	if claims.UserID != userID {
		return errors.New("token doesn't belong to this user")
	}

	// Tambahkan token ke blacklist
	blacklistedToken := &entity.Token{
		UserID:    userID,
		Token:     token,
		Type:      entity.Blacklisted,
		ExpiresAt: claims.ExpiresAt.Time,
		CreatedAt: time.Now(),
	}

	return a.tokenRepo.Create(ctx, blacklistedToken)
}

// ValidateFcmToken implements AuthUseCase.
func (a *authUseCase) ValidateFcmToken(ctx context.Context, userEmail, fcmToken string) error {
	tokenDatabase, err := a.userRepo.FindFcmTokenByEmail(ctx, userEmail)
	if err != nil {
		return err
	}

	if tokenDatabase != fcmToken {
		err := a.userRepo.UpdateFcmTokenByEmail(ctx, userEmail, fcmToken)
		if err != nil {
			return err
		}
	}

	return nil
}
