package api

import (
	"encoding/json"
	"net/http"
	"tour-service/internal/dto"
	"tour-service/internal/service" 
)

type Handler struct {
	Service *service.TourService
}

func NewHandler(s *service.TourService) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) CreateTour(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)

	var req dto.CreateTourRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tour, err := h.Service.CreateTour(userID, req)
	if err != nil {
		http.Error(w, "Failed to create tour: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tour)
}

func (h *Handler) GetMyTours(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)

	tours, err := h.Service.GetToursByAuthor(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve tours", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tours)
}