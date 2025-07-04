package handler

import (
	"backend/internal/entity"
	"backend/internal/infrastructure/oauth"
	"backend/internal/infrastructure/storage"
	"backend/internal/usecase"
	"backend/pkg/jwt"
	"mime/multipart"

	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
	googleOauth *oauth.GooogleOauth
	jwtService  *jwt.JWTService
	storage     *storage.S3Storage
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, jwtService *jwt.JWTService, googleOauth *oauth.GooogleOauth, storage *storage.S3Storage) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		googleOauth: googleOauth,
		jwtService:  jwtService,
		storage:     storage,
	}
}

type RegisterRequest struct {
	Email        string                `json:"email"`
	Password     string                `json:"password"`
	Name         string                `json:"name"`
	Role         string                `json:"role"`
	DateOfBirth  string                `json:"date_of_birth"`
	ProfilePhoto *multipart.FileHeader `json:"profile_photo"`
	PhoneNumber  string                `json:"phone_number"`
	Gender       string                `json:"gender"`
	Address      string                `json:"address"`
	BloodType    string                `json:"blood_type"`
	Rhesus       string                `json:"rhesus"`
	FCMToken     string                `json:"fcm_token"`
}

type RegisterResponse struct {
	Token string       `json:"token"`
	User  *entity.User `json:"user"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	// FCMToken string `json:"fcm_token"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *entity.User `json:"user"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type VerifyResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

type GoogleLoginRequest struct {
	Token string `json:"token"`
}

type LogoutRequest struct {
	Token string `json:"token"` // Untuk mobile/client yang perlu mengirim token
}

func isImageFile(contentType string) bool {
	// Daftar content type gambar yang diizinkan
	allowedTypes := map[string]bool{
		"image/jpeg":    true,
		"image/jpg":     true,
		"image/png":     true,
		"image/gif":     true,
		"image/webp":    true,
		"image/svg+xml": true,
		"image/bmp":     true,
		"image/tiff":    true,
	}

	return allowedTypes[contentType]
}

// Create godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags Auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param name formData string true "Name" default(Fahrul)
// @Param email formData string true "Email" default(example@example.com)
// @Param password formData string true "Password" default(fahrul123)
// @Param role formData string true "Role" default(pendonor)
// @Param date_of_birth formData string true "Date of Birth" format(date-time) default(2000-01-02)
// @Param phone_number formData string true "Phone Number" default(1234567890)
// @Param profile_photo formData file false "Profile Photo"
// @Param gender formData string true "Gender" default(male)
// @Param address formData string true "Address" default(Jakarta)
// @Param blood_type formData string false "Blood Type" default(AB)
// @Param rhesus formData string false "Rhesus" default(positive)
// @Success 201 {object} entity.User
// @Failure 400 {object} map[string]string
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// var r RegisterRequest
	// err := json.NewDecoder(r.Body).Decode(&req)

	var fileInfo *storage.FileInfo

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("profile_photo")

	if err == nil { // Hanya proses jika tidak ada error (file dikirim)
		defer file.Close()

		contentType := fileHeader.Header.Get("Content-Type")
		if !isImageFile(contentType) {
			http.Error(w, "File must be an image (JPEG, PNG, GIF, etc)", http.StatusBadRequest)
			return
		}

		fileName := fmt.Sprintf("profiles/%d_%s", time.Now().Unix(), fileHeader.Filename)
		fileInfo, err = h.storage.SaveFile(r.Context(), fileName, file, fileHeader.Header.Get("Content-Type"))
		if err != nil {
			http.Error(w, "Failed to upload file", http.StatusInternalServerError)
			return
		}
	}

	var fileURL string
	if fileInfo != nil {
		fileURL = fileInfo.URL
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("name")
	role := r.FormValue("role")
	dateOfBirth := r.FormValue("date_of_birth")
	phoneNumber := r.FormValue("phone_number")
	gender := r.FormValue("gender")
	address := r.FormValue("address")
	bloodType := r.FormValue("blood_type")
	rhesus := r.FormValue("rhesus")
	fcmToken := r.FormValue("fcm_token")

	birthDate, err := time.Parse("2006-01-02", dateOfBirth)

	if email == "" || password == "" || name == "" || role == "" || dateOfBirth == "" || phoneNumber == "" || gender == "" || address == "" {
		http.Error(w, "Please provide all required fields", http.StatusBadRequest)
		return
	}

	user, token, err := h.authUseCase.Register(r.Context(), email, password, role, name, birthDate, fileURL, phoneNumber, gender, address, bloodType, rhesus, fcmToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := RegisterResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// @Summary Login a user
// @Description Login a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email and password are required"})
		return
	}

	token, err := h.authUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal Server Error"})
		return
	}

	// err = h.authUseCase.ValidateFcmToken(r.Context(), req.Email, req.FCMToken)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Get user data
	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal Server Error"})
		return
	}

	user, err := h.authUseCase.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal Server Error"})
		return
	}

	response := LoginResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetGoogleAuthURL(w http.ResponseWriter, r *http.Request) {
	url := h.googleOauth.GetAuthURL()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}

// @Summary Google login
// @Description Google login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body GoogleLoginRequest true "Google login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Router /api/auth/google/login [post]
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// var req GoogleLoginRequest
	// req.Code = r.URL.Query().Get("code")

	// if req.Code == "" {
	// 	http.Error(w, "Code is required", http.StatusBadRequest)
	// 	return
	// }

	// token, user, err := h.authUseCase.GoogleLogin(r.Context(), req.Code)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// response := LoginResponse{
	// 	Token: token,
	// 	User:  user,
	// }

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)

	var req GoogleLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body"})
		return
	}

	userGoogle, err := h.googleOauth.VerifyGoogleIDToken(r.Context(), req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal Server Error"})
		return
	}

	token, user, err := h.authUseCase.GoogleLoginMobile(r.Context(), *userGoogle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal Server Error"})
		return
	}

	response := LoginResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
// 	var req VerifyEmailRequest
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if req.Token == "" {
// 		http.Error(w, "Token is required", http.StatusBadRequest)
// 		return
// 	}

// 	err = h.authUseCase.VerifyEmail(r.Context(), req.Token)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Email verified successfully"})
// }

// @Summary RequestPasswordReset
// @Description RequestPasswordReset
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Request Password Reset" default(example@example.com)
// @Success 200 {object} string
// @Failure 400 {object} map[string]string
// @Router /api/auth/forgot-password [post]
func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body"})
		return
	}

	if req.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email is required"})
		return
	}

	err = h.authUseCase.RequestPasswordReset(r.Context(), req.Email)
	if err != nil {
		// jangan tampilkan email ada atau tidak untuk keamanan
		// hanya tampilkan selalu sukses
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "If your email exist in our system, you will receive a password reset link shortly"})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "If your email exist in our system, you will receive a password reset link shortly"})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req VerifyResetPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" || req.NewPassword == "" {
		http.Error(w, "Token and password are required", http.StatusBadRequest)
		return
	}

	err = h.authUseCase.ResetPassword(r.Context(), req.Token, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset successfully"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Dapatkan user ID dari context middleware
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Untuk web (token di header), untuk mobile (mungkin di body)
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// Jika token kosong, cek dari body (untuk mobile)
	if token == "" {
		var req LogoutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			token = req.Token
		}
	}

	if token == "" {
		http.Error(w, "token required", http.StatusBadRequest)
		return
	}

	err := h.authUseCase.Logout(r.Context(), userID, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Clear cookie jika web
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "successfully logged out",
	})
}
