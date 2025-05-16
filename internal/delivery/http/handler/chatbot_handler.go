package handler

import (
	"backend/internal/entity"
	"backend/internal/usecase"
	"encoding/json"
	"net/http"
)

type ChatbotHandler struct {
	Usecase usecase.ChatbotUsecase
}

func NewChatbotHandler(uc usecase.ChatbotUsecase) *ChatbotHandler {
	return &ChatbotHandler{
		Usecase: uc,
	}
}

func (h *ChatbotHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	var req entity.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Format request salah"})
		return
	}

	reply, err := h.Usecase.GetReply(r.Context(), uint(req.UserID), req.Message)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	res := entity.ChatResponse{
		Reply: reply,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
