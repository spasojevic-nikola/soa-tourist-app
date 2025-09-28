package api

import (
	"encoding/json"
	"net/http"
	"tour-service/internal/models"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) CreateTour(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)

	var tour models.Tour
	if err := json.NewDecoder(r.Body).Decode(&tour); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tour.AuthorID = userID

	if err := h.DB.Create(&tour).Error; err != nil {
		http.Error(w, "Failed to create tour", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tour)
}