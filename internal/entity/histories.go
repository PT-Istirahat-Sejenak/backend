package entity

import "time"

type History struct {
	ID             uint      `json:"id"`
	UserID         uint      `json:"user_id"`
	BloodRequestID uint      `json:"blood_request_id"`
	ImageDonor     string    `json:"image_donor"`
	NextDonation   time.Time `json:"next_donation"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
