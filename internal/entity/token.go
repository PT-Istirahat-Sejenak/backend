package entity

type TokenType string

const (
	ResetPassword TokenType = "reset_password"
	EmailVerify   TokenType = "email_verify"
)

type Token struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Token     string    `json:"token"`
	Type      TokenType `json:"type"`
	ExpiredAt string    `json:"expired_at"`
	CreatedAt string    `json:"created_at"`
}
