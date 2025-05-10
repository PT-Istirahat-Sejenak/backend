package handler

import (
	"backend/internal/entity"
	"backend/internal/usecase"
	"encoding/json"
	"net/http"
)

type FcmHandler struct {
	fcmUseCase usecase.FCMUseCase
}

func NewFcmHandler(fcmUseCase usecase.FCMUseCase) *FcmHandler {
	return &FcmHandler{
		fcmUseCase: fcmUseCase,
	}
}

func (f *FcmHandler) SendFCM(w http.ResponseWriter, r *http.Request) {
	var req entity.RequestFcm
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	err = f.fcmUseCase.SendFCMV1(r.Context(), req.UserID, req.Title, req.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Success",
	})
}
