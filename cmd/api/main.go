package main

import (
	"backend/configs"
	"backend/internal/delivery/http/handler"
	"backend/internal/delivery/http/middleware"
	"backend/internal/delivery/http/routes"
	"backend/internal/infrastructure/database"
	"backend/internal/infrastructure/email"
	"backend/internal/infrastructure/oauth"
	"backend/internal/repository/postgres"
	"backend/internal/usecase"
	"backend/pkg/jwt"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		config.Google.ClientSecret,
		config.Google.RedirectURL,
	)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, emailService, googleOauth)

	// Initialize HTTP handlers
	authHandler := handler.NewAuthHandler(authUseCase, jwtService, googleOauth)
	authMiddleware := middleware.NewAuthMiddleware(jwtService, tokenRepo)

	// Initialize router
	router := mux.NewRouter()
	routes.SetupRoutes(router, authHandler, authMiddleware)

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

	log.Println("Server exited properly")
}
