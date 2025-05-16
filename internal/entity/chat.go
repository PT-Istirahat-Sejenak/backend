package entity

type ChatRequest struct {
	Message string `json:"message"`
	UserID  uint   `json:"user_id"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

type ChatMessage struct {
	Role    string `json:"role"` // "user" atau "model"
	Content string `json:"content"`
}

type ChatHistory struct {
	UserID   uint
	Messages []ChatMessage
}

