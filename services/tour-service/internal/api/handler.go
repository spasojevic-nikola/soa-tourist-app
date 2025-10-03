package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tour-service/internal/dto"
	"tour-service/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	TourService     *service.TourService
	KeyPointService *service.KeyPointService
}

func NewHandler(s *service.TourService, kps *service.KeyPointService) *Handler {
	return &Handler{
		TourService:     s,
		KeyPointService: kps,
	}
}

// CreateTour kreira turu SA keypointsima
func (h *Handler) CreateTour(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)

	var req dto.CreateTourRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tour, err := h.TourService.CreateTour(userID, req)
	if err != nil {
		http.Error(w, "Failed to create tour: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tour)
}

func (h *Handler) GetMyTours(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)

	tours, err := h.TourService.GetToursByAuthor(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve tours", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tours)
}

// vraca sve publishovane ture
func (h *Handler) GetAllPublishedTours(w http.ResponseWriter, r *http.Request) {
	tours, err := h.TourService.GetAllPublishedTours()
	if err != nil {
		http.Error(w, "Failed to retrieve published tours", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tours)
}

// KEYPOINT metode:
//
//	vraca sve key pointove za specificnu turu
func (h *Handler) GetKeyPointsByTour(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tourIDStr := vars["tourId"]
	tourID, err := strconv.ParseUint(tourIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	keyPoints, err := h.KeyPointService.GetKeyPointsByTour(uint(tourID))
	if err != nil {
		http.Error(w, "Failed to retrieve key points: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keyPoints)
}

// azurira keypoint
func (h *Handler) UpdateKeyPoint(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)

	vars := mux.Vars(r)
	keyPointIDStr := vars["keyPointId"]
	keyPointID, err := strconv.ParseUint(keyPointIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid key point ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateKeyPointRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	keyPoint, err := h.KeyPointService.UpdateKeyPoint(uint(keyPointID), userID, req)
	if err != nil {
		http.Error(w, "Failed to update key point: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keyPoint)
}

// brise keypoint
func (h *Handler) DeleteKeyPoint(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)

	vars := mux.Vars(r)
	keyPointIDStr := vars["keyPointId"]
	keyPointID, err := strconv.ParseUint(keyPointIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid key point ID", http.StatusBadRequest)
		return
	}

	err = h.KeyPointService.DeleteKeyPoint(uint(keyPointID), userID)
	if err != nil {
		http.Error(w, "Failed to delete key point: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// dodaje duration u turu
func (h *Handler) AddDuration(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)
	vars := mux.Vars(r)
	tourIDStr := vars["tourId"]
	tourID, err := strconv.ParseUint(tourIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	var req dto.AddDurationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	duration, err := h.TourService.AddDuration(uint(tourID), userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// posl dodavanja trajanja, rekalkulise distancu ako ima 2+ keypointa
	h.TourService.CalculateDistance(uint(tourID))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(duration)
}

// publishuje draftovanu turu
func (h *Handler) PublishTour(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)
	vars := mux.Vars(r)
	tourIDStr := vars["tourId"]
	tourID, err := strconv.ParseUint(tourIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	tour, err := h.TourService.PublishTour(uint(tourID), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tour)
}

// arhivira publishovanu turu
func (h *Handler) ArchiveTour(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)
	vars := mux.Vars(r)
	tourIDStr := vars["tourId"]
	tourID, err := strconv.ParseUint(tourIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	tour, err := h.TourService.ArchiveTour(uint(tourID), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tour)
}

// reaktivira arhiviranu turu
func (h *Handler) ActivateTour(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)
	vars := mux.Vars(r)
	tourIDStr := vars["tourId"]
	tourID, err := strconv.ParseUint(tourIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	tour, err := h.TourService.ActivateTour(uint(tourID), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tour)
}

// vraca turu po id sa svim keypointima i trajanjem
func (h *Handler) GetTourByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tourIDStr := vars["tourId"]
	tourID, err := strconv.ParseUint(tourIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	tour, err := h.TourService.GetTourByID(uint(tourID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tour)
}
