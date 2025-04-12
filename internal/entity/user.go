package entity

import "time"

type User struct {
	ID              uint      `json:"id"`
	Email           string    `json:"email"`
	Password        string    `json:"-"`
	Name            string    `json:"name"`
	GoogleID        *string   `json:"-"`
	IsEmailVerified bool      `json:"is_email_verified"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewGoogleUser(email, name, googleID string, verified bool) *User {
	return &User{
		Email:           email,
		Name:            name,
		GoogleID:        &googleID,
		IsEmailVerified: verified,
	}
}
