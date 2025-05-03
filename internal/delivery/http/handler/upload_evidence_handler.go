package handler

import (
	"backend/internal/infrastructure/storage"
	"backend/internal/usecase"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type UploadEvidenceHandler struct {
	uploadEvidenceUseCase usecase.EvidenceUseCase
	storage               *storage.S3Storage
}

func NewUploadEvidenceHandler(uploadEvidenceUseCase usecase.EvidenceUseCase, storage *storage.S3Storage) *UploadEvidenceHandler {
	return &UploadEvidenceHandler{
		uploadEvidenceUseCase: uploadEvidenceUseCase,
		storage:               storage,
	}
}

type UploadEvidenceRequest struct {
	File   *multipart.FileHeader `json:"image"`
	UserID uint                  `json:"user_id"`
}

// @Summary Upload Evidence
// @Description Upload Evidence
// @Tags Upload Evidence
// @Accept x-www-form-urlencoded
// @Produce json
// @Param image formData file true "Image"
// @Param user_id formData string true "User ID" default(1)
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/upload-evidence [post]
func (h *UploadEvidenceHandler) PostUploadEvidence(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := fmt.Sprintf("evidences/%d_%s", time.Now().Unix(), fileHeader.Filename)
	fileInfo, err := h.storage.SaveFile(r.Context(), fileName, file, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	// Ambil user ID dari context (diset oleh middleware auth)
	// userID, ok := r.Context().Value("user_id").(uint)
	// if !ok {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	userID, err := strconv.ParseUint(r.FormValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.uploadEvidenceUseCase.UploadEvidence(r.Context(), uint(userID), fileInfo.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded successfully"})
}
