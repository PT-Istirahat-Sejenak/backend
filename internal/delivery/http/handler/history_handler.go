package handler

import (
	"backend/internal/infrastructure/storage"
	"backend/internal/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type HistoryHandler struct {
	historyUseCase usecase.HistoryUseCase
	authUseCase    usecase.AuthUseCase
	storage        *storage.S3Storage
}

func NewHistoryHandler(historyUseCase usecase.HistoryUseCase, authUseCase usecase.AuthUseCase, storage *storage.S3Storage) *HistoryHandler {
	return &HistoryHandler{
		historyUseCase: historyUseCase,
		authUseCase:    authUseCase,
		storage:        storage,
	}
}

type HistoryRequest struct {
	UserID         uint   `json:"user_id"`
	BloodRequestID uint   `json:"blood_request_id"`
	ImageDonor     string `json:"image_donor"`
}

type HistoryRequestReponse struct {
	ID             uint      `json:"id"`
	BloodRequestID uint      `json:"blood_request_id"`
	ImageDonor     string    `json:"image_donor"`
	NextDonation   time.Time `json:"next_donation"`
}

type TimeHistoryRequest struct {
	UserID uint `json:"user_id"`
}

type NextDonationResponse struct {
	NextDonation time.Time `json:"next_donation"`
}

type LatestDonationResponse struct {
	LatestDonation time.Time `json:"latest_donation"`
}

type HistoryResponse struct {
	ID             uint      `json:"id"`
	UserID         uint      `json:"user_id"`
	BloodRequestID uint      `json:"blood_request_id"`
	ImageDonor     string    `json:"image_donor"`
	NextDonation   time.Time `json:"next_donation"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// @Summary Create History
// @Description Create a new History
// @Tags History
// @Accept multipart/form-data
// @Produce json
// @Param image_donor formData file true "Image Donor"
// @Param user_id formData string true "User ID"
// @Param blood_request_id formData string true "Blood Request ID"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/history [post]
func (h *HistoryHandler) PostHistory(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body"})
		return
	}

	file, fileHeader, err := r.FormFile("image_donor")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to get file from form"})
		return
	}
	defer file.Close()

	// if file != nil {
	// 	fileName := fmt.Sprintf("histories/%d_%s", time.Now().Unix(), fileHeader.Filename)
	// 	fileInfo, err := h.storage.SaveFile(r.Context(), fileName, file, fileHeader.Header.Get("Content-Type"))
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to upload file"})
	// 		return
	// 	}

	// 	var historyRequest HistoryRequest
	// 	err = json.NewDecoder(r.Body).Decode(&historyRequest)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body"})
	// 		return
	// 	}

	// 	userID, err := strconv.ParseUint(r.FormValue("user_id"), 10, 64)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid user ID"})
	// 		return
	// 	}
	// 	bloodRequestID, err := strconv.ParseUint(r.FormValue("blood_request_id"), 10, 64)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid blood request ID"})
	// 		return
	// 	}

	// 	now := time.Now()
	// 	user, err := h.authUseCase.GetUserByID(r.Context(), uint(userID))
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		json.NewEncoder(w).Encode(map[string]string{"message": "Internal Server Error"})
	// 		return
	// 	}

	// 	var nextDonation time.Time
	// 	if user.Gender == "male" {
	// 		nextDonation = now.AddDate(0, 3, 0)
	// 	} else if user.Gender == "female" {
	// 		nextDonation = now.AddDate(0, 4, 0)
	// 	} else {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid Data"})
	// 		return
	// 	}

	// 	err = h.historyUseCase.AddHistory(r.Context(), uint(userID), uint(bloodRequestID), fileInfo.URL, nextDonation)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to create history"})
	// 		return
	// 	}

	// 	w.WriteHeader(http.StatusCreated)
	// 	json.NewEncoder(w).Encode(map[string]string{"message": "History created"})
	// 	return
	// }

	// var historyRequest HistoryRequest
	// err = json.NewDecoder(r.Body).Decode(&historyRequest)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body"})
	// 	return
	// }

	// userID, err := strconv.ParseUint(r.FormValue("user_id"), 10, 64)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(map[string]string{"message": "Invalid user ID"})
	// 	return
	// }
	// bloodRequestID, err := strconv.ParseUint(r.FormValue("blood_request_id"), 10, 64)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(map[string]string{"message": "Invalid blood request ID"})
	// 	return
	// }

	// now := time.Now()
	// user, err := h.authUseCase.GetUserByID(r.Context(), uint(userID))
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(map[string]string{"message": "Internal Server Error"})
	// 	return
	// }

	// var nextDonation time.Time
	// if user.Gender == "male" {
	// 	nextDonation = now.AddDate(0, 3, 0)
	// } else if user.Gender == "female" {
	// 	nextDonation = now.AddDate(0, 4, 0)
	// } else {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(map[string]string{"message": "Invalid Data"})
	// 	return
	// }

	// err = h.historyUseCase.AddHistory(r.Context(), uint(userID), uint(bloodRequestID), "", nextDonation)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(map[string]string{"message": "Failed to create history"})
	// 	return
	// }

	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(map[string]string{"message": "History created"})

	contentType := fileHeader.Header.Get("Content-Type")
	if !isImageFile(contentType) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid file type"})
		return
	}
	fileName := fmt.Sprintf("histories/%d_%s", time.Now().Unix(), fileHeader.Filename)
	fileInfo, err := h.storage.SaveFile(r.Context(), fileName, file, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to upload file"})
		return
	}

	var historyRequest HistoryRequest
	err = json.NewDecoder(r.Body).Decode(&historyRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body"})
		return
	}

	userID, err := strconv.ParseUint(r.FormValue("user_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid user ID"})
		return
	}
	bloodRequestID, err := strconv.ParseUint(r.FormValue("blood_request_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid blood request ID"})
		return
	}

	now := time.Now()
	user, err := h.authUseCase.GetUserByID(r.Context(), uint(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal Server Error"})
		return
	}

	var nextDonation time.Time
	if user.Gender == "male" {
		nextDonation = now.AddDate(0, 3, 0)
	} else if user.Gender == "female" {
		nextDonation = now.AddDate(0, 4, 0)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid Data"})
		return
	}

	err = h.historyUseCase.AddHistory(r.Context(), uint(userID), uint(bloodRequestID), fileInfo.URL, nextDonation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to create history"})
		return
	}

	total := user.TotalDonation + 1
	coin := user.Coin + 10
	err = h.authUseCase.UpdateCountDonation(r.Context(), uint(userID), total)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to update count donation"})
	}

	err = h.authUseCase.UpdateCoinTotal(r.Context(), uint(userID), coin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to update coin"})
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "History created"})
}

// @Summary Get user history
// @Description Retrieves a list of history records for a given user ID.
// @Tags History
// @Accept json
// @Produce json
// @Param user_id body uint true "User ID"
// @Success 200 {array} HistoryResponse
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/history [get]
func (h *HistoryHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	var req TimeHistoryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid user ID"})
		return
	}

	histories, err := h.historyUseCase.HistoryByUserId(r.Context(), uint(req.UserID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal Server Error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(histories)
}

// @Summary Get latest history
// @Description Retrieves the latest history record for a given user ID.
// @Tags History
// @Accept json
// @Produce json
// @Param user_id body uint true "User ID"
// @Success 200 {object} time.Time "Latest donation date"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/history/latest [get]
func (h *HistoryHandler) GetLatestHistory(w http.ResponseWriter, r *http.Request) {
	var userID uint
	err := json.NewDecoder(r.Body).Decode(&userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid user ID"})
		return
	}

	history, err := h.historyUseCase.GetLatestDonation(r.Context(), uint(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal Server Error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}

// @Summary Get next donation date
// @Description Retrieves the next scheduled donation date for the specified user.
// @Tags History
// @Accept json
// @Produce json
// @Param user_id query uint true "User ID"
// @Success 200 {object} time.Time "Next donation date"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/history/next [get]
func (h *HistoryHandler) GetNextHistory(w http.ResponseWriter, r *http.Request) {
	var userID uint
	err := json.NewDecoder(r.Body).Decode(&userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid user ID"})
		return
	}

	history, err := h.historyUseCase.GetNextDonation(r.Context(), uint(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal Server Error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}
