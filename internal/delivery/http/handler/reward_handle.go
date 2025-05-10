package handler

import (
	"backend/internal/usecase"
	"encoding/json"
	"net/http"
)

type RewardHandler struct {
	Usecase usecase.RewardUseCase
}

type RewardRequestt struct {
	UserID     uint   `json:"user_id"`
	OperatorID string `json:"operator_id"`
	Amount     string `json:"amount"`
	Number     string `json:"number"`
}

func NewRewardHandler(rewardUseCase usecase.RewardUseCase) *RewardHandler {
	return &RewardHandler{
		Usecase: rewardUseCase,
	}
}

func (rw *RewardHandler) ClaimReward(w http.ResponseWriter, r *http.Request) {
	var req RewardRequestt
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	res, err := rw.Usecase.GetReward(r.Context(), req.UserID, req.OperatorID, req.Amount, req.Number)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (rw *RewardHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	res, err := rw.Usecase.GetBalance(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
