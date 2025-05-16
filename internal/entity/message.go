package entity

import "time"

type Message struct {
	ID          uint      `json:"id"`
	SenderID    uint      `json:"sender_id"`
	ReceiverID  uint      `json:"receiver_id"`
	Content     string    `json:"content"`
	IsDelivered bool      `json:"is_delivered"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MessageRequest struct {
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type ChatHistoryRequest struct {
	UserID1 uint `json:"user_id_1"`
	UserID2 uint `json:"user_id_2"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
}
