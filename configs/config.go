package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Email    EmailConfig
	Google   GoogleOAuthConfig
	Storage  StorageConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	ExpireTime int
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SenderEmail  string
	SenderName   string
	SMTPPassword string
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type StorageConfig struct {
	Type          string // "local" atau "s3"
	LocalBasePath string
	LocalBaseURL  string
	S3AccountID   string
	S3AccessKey   string
	S3SecretKey   string
	S3BucketName  string
	S3Region      string
	S3BaseURL     string
}

func LoadConfig() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("error loading .env file: %v", err)
		}
	}

	expTime, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_TIME"))

	return &Config{
		Server: ServerConfig{
			Port: os.Getenv("PORT"),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSL_MODE"),
		},
		JWT: JWTConfig{
			Secret:     os.Getenv("JWT_SECRET"),
			ExpireTime: expTime,
		},
		Email: EmailConfig{
			SMTPHost:     os.Getenv("SMTP_HOST"),
			SMTPPort:     os.Getenv("SMTP_PORT"),
			SenderEmail:  os.Getenv("SENDER_EMAIL"),
			SenderName:   os.Getenv("SENDER_NAME"),
			SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		},
		Google: GoogleOAuthConfig{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		},
		Storage: StorageConfig{
			Type:          os.Getenv("STORAGE_TYPE"),
			LocalBasePath: os.Getenv("LOCAL_STORAGE_PATH"),
			LocalBaseURL:  os.Getenv("LOCAL_STORAGE_URL"),
			S3AccountID:   os.Getenv("S3_ACCOUNT_ID"),
			S3AccessKey:   os.Getenv("S3_ACCESS_KEY"),
			S3SecretKey:   os.Getenv("S3_SECRET_KEY"),
			S3BucketName:  os.Getenv("S3_BUCKET_NAME"),
			S3Region:      os.Getenv("S3_REGION"),
			S3BaseURL:     os.Getenv("S3_BASE_URL"),
		},
	}, nil
}
