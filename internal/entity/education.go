package entity

import "time"

type Education struct {
	ID        uint      `json:"id"`
	Image     string    `json:"image"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
