package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"backend/configs"
	_ "backend/docs"
	"backend/internal/delivery/http/handler"
	"backend/internal/delivery/http/middleware"
	"backend/internal/delivery/http/routes"
	"backend/internal/infrastructure/database"
	"backend/internal/infrastructure/email"
	"backend/internal/infrastructure/oauth"
	"backend/internal/infrastructure/storage"
	"backend/internal/repository/postgres"
	"backend/internal/usecase"
	"backend/pkg/jwt"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database connection
	db, err := database.NewPostgresConnection(&config.Database)
	if err != nil {
		log.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize database schema
	err = database.InitDatabase(db)
	if err != nil {
		log.Fatal("Failed to initialize database schema: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)
	educationRepo := postgres.NewEducationRepository(db)
	uploadEvidenceRepo := postgres.NewUploadEvidenceRepository(db)
	historyRepo := postgres.NewHistoryRepository(db)
	rewardRepo := postgres.NewRewardRepository(db)
	messageRepo := postgres.NewMessageRepository(db)

	// Initialize services
	jwtService := jwt.NewJWTService(config.JWT.Secret, config.JWT.ExpireTime)
	emailService := email.NewEmialService(
		config.Email.SMTPHost,
		config.Email.SMTPPort,
		config.Email.SenderEmail,
		config.Email.SenderName,
		config.Email.SMTPPassword,
	)
	googleOauth := oauth.NewGoogleOauth(
		config.Google.ClientID,
		// config.Google.ClientSecret,
		config.Google.RedirectURL,
	)

	// Initialize storage

	var fileStorage storage.FileStorage

	if config.Storage.Type == "s3" {
		fileStorage, err = storage.NewS3Storage(
			config.Storage.S3AccountID,
			config.Storage.S3AccessKey,
			config.Storage.S3SecretKey,
			config.Storage.S3BucketName,
			config.Storage.S3Region,
			config.Storage.S3BaseURL,
		)
	} else {
		// Default to local storage
		fileStorage, err = storage.NewLocalStorage(
			config.Storage.LocalBasePath,
			config.Storage.LocalBaseURL,
		)
	}

	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize use cases
	fcmUseCase := usecase.NewFcmUseCase(userRepo)
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, fcmUseCase, emailService, googleOauth)
	profileUseCase := usecase.NewProfileUseCase(userRepo, fileStorage)
	educationUseCase := usecase.NewEducationUseCase(educationRepo, fileStorage)
	uploadEvidenceUseCase := usecase.NewUploadEvidenceUseCase(uploadEvidenceRepo, fileStorage)
	historyUseCase := usecase.NewHistoryUseCase(historyRepo, fileStorage)
	chatbotUseCase := usecase.NewChatbotUsecase(config.ChatBot)
	rewardUseCase := usecase.NewRewardUseCase(config.Reloadly, userRepo, rewardRepo)
	messageUseCase := usecase.NewMessageUseCase(messageRepo)

	// Initialize HTTP handlers
	authHandler := handler.NewAuthHandler(authUseCase, jwtService, googleOauth, fileStorage.(*storage.S3Storage))
	authMiddleware := middleware.NewAuthMiddleware(jwtService, tokenRepo)
	profileHandler := handler.NewProfileHandler(profileUseCase)
	educationHandler := handler.NewEducationHandler(educationUseCase, fileStorage.(*storage.S3Storage))
	uploadEvidenceHandler := handler.NewUploadEvidenceHandler(uploadEvidenceUseCase, fileStorage.(*storage.S3Storage))
	historyHandler := handler.NewHistoryHandler(historyUseCase, authUseCase, fileStorage.(*storage.S3Storage))
	chatbotHandler := handler.NewChatbotHandler(chatbotUseCase)
	rewardHanlder := handler.NewRewardHandler(rewardUseCase)
	fcmHandler := handler.NewFcmHandler(fcmUseCase)
	messageHandler := handler.NewWebSockerHandler(messageUseCase)

	go messageHandler.HandleMessages()

	// Initialize router
	router := mux.NewRouter()
	routes.SetupRoutes(router, authHandler, authMiddleware, profileHandler, educationHandler, uploadEvidenceHandler, historyHandler, chatbotHandler, rewardHanlder, fcmHandler, messageHandler)

	// Configure HTTP server
	server := &http.Server{
		Addr:         ":" + config.Server.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start HTTP server
	go func() {
		log.Printf("Starting server on port %s", config.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Doesn't block if no connection, but will otherwise wait until the timeout
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Jika menggunakan local storage, tambahkan route untuk static files
	if config.Storage.Type == "local" {
		// Buat direktori uploads jika belum ada
		uploadsDir := filepath.Join(".", "uploads")
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			log.Fatalf("Failed to create uploads directory: %v", err)
		}

		// Serve static files dari direktori uploads
		fs := http.FileServer(http.Dir(uploadsDir))
		router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fs))
	}

	log.Println("Server exited properly")
}
