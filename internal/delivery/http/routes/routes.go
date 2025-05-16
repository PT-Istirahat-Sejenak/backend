package routes

import (
	"backend/internal/delivery/http/handler"
	"backend/internal/delivery/http/middleware"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(
	router *mux.Router,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
	profileHandler *handler.ProfileHandler,
	eduHandler *handler.EducationHandler,
	edivenceHandler *handler.UploadEvidenceHandler,
	historyHandler *handler.HistoryHandler,
	chatBotHandler *handler.ChatbotHandler,
	rewardHandler *handler.RewardHandler,
	fcmHandler *handler.FcmHandler,
	websockerHandler *handler.WebSocketHandler,

) {
	// Public routes
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/api/auth/google/url", authHandler.GetGoogleAuthURL).Methods("GET")
	router.HandleFunc("/api/auth/google/callback", authHandler.GoogleLogin).Methods("GET")
	// router.HandleFunc("/api/auth/verify-email", authHandler.VerifyEmail).Methods("POST")
	router.HandleFunc("/api/auth/forgot-password", authHandler.RequestPasswordReset).Methods("POST")
	router.HandleFunc("/api/auth/reset-password", authHandler.ResetPassword).Methods("POST")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// cek education without middleware
	router.HandleFunc("/api/educations", eduHandler.GetEducations).Methods("GET")
	router.HandleFunc("/api/education", eduHandler.PostEducation).Methods("POST")
	router.HandleFunc("/api/education", eduHandler.Update).Methods("PUT")
	router.HandleFunc("/api/education", eduHandler.Delete).Methods("Delete")

	// cek upload evidence image
	router.HandleFunc("/api/upload-evidence", edivenceHandler.PostUploadEvidence).Methods("POST")

	// cek history tanpa middleware
	router.HandleFunc("/api/history", historyHandler.PostHistory).Methods("POST")

	// cek history tanpa middleware
	router.HandleFunc("/api/history", historyHandler.GetHistory).Methods("GET")
	router.HandleFunc("/api/history", historyHandler.PostHistory).Methods("POST")
	router.HandleFunc("/api/history/latest", historyHandler.GetLatestHistory).Methods("GET")
	router.HandleFunc("/api/history/next", historyHandler.GetNextHistory).Methods("GET")

	// cek chatbot tanpa middleware
	router.HandleFunc("/api/chatbot", chatBotHandler.HandleChat).Methods("POST")

	// Reward tanpa middlewar
	router.HandleFunc("/api/reward", rewardHandler.ClaimReward).Methods("POST")
	router.HandleFunc("/api/reward/balance", rewardHandler.GetBalance).Methods("GET")

	// Fcm tanpa middleware
	router.HandleFunc("/api/broadcast", fcmHandler.SendFCM).Methods("POST")

	// Message tanpa middleware
	router.HandleFunc("/api/message", websockerHandler.HandleConnection)

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

	// education routes
	// protected.HandleFunc("/educations",eduHandler.GetEducations).Methods("GET")
	// protected.HandleFunc("/api/educations-pedonor", eduHandler.GetEducationsPendonor).Methods("GET")
	// protected.HandleFunc("/api/educations-pencari-donor", eduHandler.GetEducationsPencariDonor).Methods("GET")
	// protected.HandleFunc("/education", eduHandler.PostEducation).Methods("POST")
	// protected.HandleFunc("/education", eduHandler.Update).Methods("PUT")
	// protected.HandleFunc("/education", eduHandler.Delete).Methods("Delete")

	// upload evidence routes
	// protected.HandleFunc("/upload-evidence", edivenceHandler.PostUploadEvidence).Methods("POST")

	// history routes
	// protected.HandleFunc("/history", historyHandler.PostHistory).Methods("POST")
	// protected.HandleFunc("/history", historyHandler.GetHistory).Methods("GET")
	// protected.HandleFunc("/history/latest", historyHandler.GetLatestHistory).Methods("GET")
	// protected.HandleFunc("/history/next", historyHandler.GetNextHistory).Methods("GET")

	// chatbot routes
	// protected.HandleFunc("/chatbot", chatBotHandler.HandleChat).Methods("POST")

	// reward routes
	// protected.HandleFunc("/reward", rewardHandler.ClaimReward).Methods("POST")
	// protected.HandleFunc("/reward/balance", rewardHandler.GetBalance).Methods("GET")

	// fcm routes
	// protected.HandleFunc("/broadcast", fcmHandler.SendFCM).Methods("POST")

	// message routes
	// protected.HandleFunc("/message", websockerHandler.HandleConnection).Methods("GET")
}
