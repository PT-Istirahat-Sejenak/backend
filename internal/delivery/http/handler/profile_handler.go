package handler

import (
	"backend/internal/usecase"
	"encoding/json"
	"net/http"
)

type ProfileHandler struct {
	profileUseCase usecase.ProfileUseCase
}

func NewProfileHandler(profileUseCase usecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{
		profileUseCase: profileUseCase,
	}
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Ambil user ID dari context (diset oleh middleware auth)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.profileUseCase.GetUserProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *ProfileHandler) UpdateProfilePhoto(w http.ResponseWriter, r *http.Request) {
	// Ambil user ID dari context (diset oleh middleware auth)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Maksimum ukuran file 5MB
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("profile_photo")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	photoURL, err := h.profileUseCase.UpdateProfilePhoto(r.Context(), userID, file, header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile photo updated successfully",
		"url":     photoURL,
	})
}
