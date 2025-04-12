package middleware

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"backend/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	jwtService *jwt.JWTService
	tokenRepo  repository.TokenRepository
}

func NewAuthMiddleware(jwtService *jwt.JWTService, tokenRepo repository.TokenRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		tokenRepo:  tokenRepo,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := bearerToken[1]
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Tambahkan pengecekan blacklist
		isBlacklisted, err := m.tokenRepo.Exists(r.Context(), token, entity.Blacklisted)
		if err != nil {
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			return
		}
		if isBlacklisted {
			http.Error(w, "token has been revoked", http.StatusUnauthorized)
			return
		}

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
