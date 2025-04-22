package routes

import (
	"backend/internal/delivery/http/handler"
	"backend/internal/delivery/http/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(
	router *mux.Router,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
	profileHandler *handler.ProfileHandler,
) {
	// Public routes
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/api/auth/google/url", authHandler.GetGoogleAuthURL).Methods("GET")
	router.HandleFunc("/api/auth/google/login", authHandler.GoogleLogin).Methods("POST")
	router.HandleFunc("/api/auth/verify-email", authHandler.VerifyEmail).Methods("POST")
	router.HandleFunc("/api/auth/forgot-password", authHandler.RequestPasswordReset).Methods("POST")
	router.HandleFunc("/api/auth/reset-password", authHandler.ResetPassword).Methods("POST")

	// Protected routes
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(authMiddleware.Authenticate)
	protected.HandleFunc("/user/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Protected route"))
	}).Methods("GET")
	protected.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")

	// profile routes
	protected.HandleFunc("/user/profile", profileHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/user/profile/photo", profileHandler.UpdateProfilePhoto).Methods("POST")
}
