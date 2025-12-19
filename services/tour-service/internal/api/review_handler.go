package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tour-service/internal/dto"
	"tour-service/internal/service"

	"github.com/gorilla/mux"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// CreateReview handles creating a new review
func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	// Get tourist ID from context (set by AuthMiddleware)
	touristID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	review, err := h.reviewService.CreateReview(touristID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

// GetReviewsByTour handles retrieving all reviews for a tour
func (h *ReviewHandler) GetReviewsByTour(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tourID, err := strconv.ParseUint(vars["tourId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	reviews, err := h.reviewService.GetReviewsByTourID(uint(tourID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// GetTourRatingStats handles retrieving rating statistics for a tour
func (h *ReviewHandler) GetTourRatingStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tourID, err := strconv.ParseUint(vars["tourId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	stats, err := h.reviewService.GetTourRatingStats(uint(tourID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// UpdateReview handles updating a review
func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	touristID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	reviewID, err := strconv.ParseUint(vars["reviewId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	review, err := h.reviewService.UpdateReview(uint(reviewID), touristID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// DeleteReview handles deleting a review
func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	touristID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	reviewID, err := strconv.ParseUint(vars["reviewId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	err = h.reviewService.DeleteReview(uint(reviewID), touristID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetMyReviews handles retrieving all reviews by the authenticated tourist
func (h *ReviewHandler) GetMyReviews(w http.ResponseWriter, r *http.Request) {
	touristID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	reviews, err := h.reviewService.GetReviewsByTouristID(touristID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}
