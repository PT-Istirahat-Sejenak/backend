package entity

import "time"

type UploadEvidence struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}