package entity

import "time"

type TokenType string

const (
	ResetPassword TokenType = "reset_password"
	EmailVerify   TokenType = "email_verify"
	Blacklisted   TokenType = "black_list"
)

type Token struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Token     string    `json:"token"`
	Type      TokenType `json:"type"`
	ExpiresAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}
