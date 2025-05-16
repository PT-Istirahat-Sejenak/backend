package entity

import "time"

type User struct {
	ID            uint      `json:"id"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	Role          string    `json:"role"`
	Name          string    `json:"name"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	ProfilePhoto  *string    `json:"profile_photo"`
	PhoneNumber   string    `json:"phone_number"`
	Gender        string    `json:"gender"`
	Address       string    `json:"address"`
	BloodType     string    `json:"blood_type"`
	Rhesus        string    `json:"rhesus"`
	GoogleID      *string   `json:"-"`
	TotalDonation int       `json:"total_donation"`
	Coin          int       `json:"coin"`
	FCMToken      string    `json:"fcm_token"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewGoogleUser(email, name, googleID string) *User {
	return &User{
		Email:    email,
		Name:     name,
		GoogleID: &googleID,
	}
}
