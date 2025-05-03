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

type EducationHandler struct {
	educationUseCase usecase.EducationUseCase
	storage          *storage.S3Storage
}

func NewEducationHandler(educationUseCase usecase.EducationUseCase, storage *storage.S3Storage) *EducationHandler {
	return &EducationHandler{
		educationUseCase: educationUseCase,
		storage:          storage,
	}
}

type EducationRequest struct {
	Image   *multipart.FileHeader `json:"image"`
	Title   string                `json:"title"`
	Content string                `json:"content"`
	Type    string                `json:"type"`
}

type EducationResponse struct {
	ID        uint      `json:"id"`
	Image     string    `json:"image"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EducationDeleteReuest struct {
	ID uint `json:"id"`
}

type EducationUpdateRequest struct {
	ID        uint                  `json:"id"`
	Image     *multipart.FileHeader `json:"image"`
	Title     string                `json:"title"`
	Content   string                `json:"content"`
	Type      string                `json:"type"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// @Summary Create a new Education
// @Description Post a new Education
// @Tags Education
// @Accept x-www-form-urlencoded
// @Produce json
// @Param image formData file true "Image"
// @Param title formData string true "Title" default(Donora menang juara 1 di GSC tingkat international 2025)
// @Param content formData string true "Content" default(This is)
// @Param type formData string true "Type" default(pendonor)
// @Success 201 {object} EducationResponse
// @Failure 400 {object} map[string]string
// @Router /api/education [post]
func (e *EducationHandler) PostEducation(w http.ResponseWriter, r *http.Request) {
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

	fileName := fmt.Sprintf("educations/%d_%s", time.Now().Unix(), fileHeader.Filename)
	fileInfo, err := e.storage.SaveFile(r.Context(), fileName, file, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	typee := r.FormValue("type")

	err = e.educationUseCase.Post(r.Context(), fileInfo.URL, title, content, typee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Education created successfully"})
}

// @Summary Get Education
// @Description Get All and get by id and type with params
// @Tags Education
// @Accept json
// @Produce json
// @Success 200 {object} EducationResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/educations [get]
func (e *EducationHandler) GetEducations(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		educations, err := e.educationUseCase.GetAllEducation(r.Context())
		if err != nil {
			http.Error(w, "Canot get education", http.StatusInternalServerError)
			return
		}

		if educations == nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Education Is Empty"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(educations)
		return
	}

	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	edu, err := e.educationUseCase.FindById(r.Context(), uint(idInt))
	if err != nil {
		http.Error(w, "Cannot find education", http.StatusInternalServerError)
		return
	}

	if edu == nil {
		http.Error(w, "Cannot find education", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(edu)
}

// @Summary Get Education Pendonor
// @Description Get Education Pendonor
// @Tags Education
// @Accept json
// @Produce json
// @Success 200 {object} EducationResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/educations-pendonor [get]
func (e *EducationHandler) GetEducationsPendonor(w http.ResponseWriter, r *http.Request) {
	educations, err := e.educationUseCase.FindEducationPendonor(r.Context())
	if err != nil {
		http.Error(w, "Canot get education", http.StatusInternalServerError)
		return
	}

	if educations == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Education Is Empty"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(educations)
}

// @Summary Get Education Pencari Donor
// @Description Get Education Pencari Donor
// @Tags Education
// @Accept json
// @Produce json
// @Success 200 {object} EducationResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/educations/pencari-donor [get]
func (e *EducationHandler) GetEducationsPencariDonor(w http.ResponseWriter, r *http.Request) {
	educations, err := e.educationUseCase.FindEducationPencariDonor(r.Context())
	if err != nil {
		http.Error(w, "Canot get education", http.StatusInternalServerError)
		return
	}

	if educations == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Education Is Empty"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(educations)
}

// @Summary Update education
// @Description Update education
// @Tags Education
// @Accept x-www-form-urlencoded
// @Produce json
// @Param title formData string true "Title" default(Donora menang juara 1 di GSC tingkat international 2025)
// @Param image formData file true "Image"
// @Param content formData string true "Content" default(Alhamdulillah)
// @Param type formData string true "Type" default(patient)
// @Success 200 {object} EducationResponse
// @Failure 400 {object} map[string]string
// @Router /api/education [put]
func (e *EducationHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	fileName := fmt.Sprintf("educations/%d_%s", time.Now().Unix(), fileHeader.Filename)
	fileInfo, err := e.storage.SaveFile(r.Context(), fileName, file, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	typee := r.FormValue("type")

	err = e.educationUseCase.Post(r.Context(), fileInfo.URL, title, content, typee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Education update successfully"})
}

// @Summary Delete education
// @Description Delete education
// @Tags Education
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/education [delete]
func (e *EducationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = e.educationUseCase.Delete(r.Context(), uint(idInt))

	if err != nil {
		http.Error(w, "Cannot delete education", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "Education deleted"})
}
