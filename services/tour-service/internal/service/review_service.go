package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"tour-service/internal/dto"
	"tour-service/internal/models"
	"tour-service/internal/repository"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	tourRepo   *repository.TourRepository
}

func NewReviewService(reviewRepo *repository.ReviewRepository, tourRepo *repository.TourRepository) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		tourRepo:   tourRepo,
	}
}

// CreateReview creates a new review
func (s *ReviewService) CreateReview(touristID uint, req dto.CreateReviewRequest) (*models.Review, error) {
	// Verify tour exists and is published
	tour, err := s.tourRepo.FindByID(req.TourID)
	if err != nil {
		return nil, errors.New("tour not found")
	}
	if tour.Status != models.Published {
		return nil, errors.New("can only review published tours")
	}

	// Get tourist username from stakeholders service
	username, err := s.getUserUsername(touristID)
	if err != nil {
		// If we can't get username, use a default
		username = fmt.Sprintf("User_%d", touristID)
	}

	// Convert images array to JSON string
	imagesJSON, err := json.Marshal(req.Images)
	if err != nil {
		return nil, err
	}

	review := &models.Review{
		TourID:          req.TourID,
		TouristID:       touristID,
		TouristUsername: username,
		Rating:          req.Rating,
		Comment:         req.Comment,
		VisitDate:       req.VisitDate,
		Images:          string(imagesJSON),
	}

	err = s.reviewRepo.Create(review)
	if err != nil {
		return nil, err
	}

	return review, nil
}

// GetReviewByID retrieves a review by ID
func (s *ReviewService) GetReviewByID(id uint) (*models.Review, error) {
	return s.reviewRepo.FindByID(id)
}

// GetReviewsByTourID retrieves all reviews for a tour
func (s *ReviewService) GetReviewsByTourID(tourID uint) ([]models.Review, error) {
	return s.reviewRepo.FindByTourID(tourID)
}

// GetReviewsByTouristID retrieves all reviews by a tourist
func (s *ReviewService) GetReviewsByTouristID(touristID uint) ([]models.Review, error) {
	return s.reviewRepo.FindByTouristID(touristID)
}

// UpdateReview updates an existing review
func (s *ReviewService) UpdateReview(reviewID, touristID uint, req dto.UpdateReviewRequest) (*models.Review, error) {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return nil, errors.New("review not found")
	}

	// Verify ownership
	if review.TouristID != touristID {
		return nil, errors.New("unauthorized to update this review")
	}

	// Update fields if provided
	if req.Rating > 0 {
		review.Rating = req.Rating
	}
	if req.Comment != "" {
		review.Comment = req.Comment
	}
	if !req.VisitDate.IsZero() {
		review.VisitDate = req.VisitDate
	}
	if req.Images != nil {
		imagesJSON, err := json.Marshal(req.Images)
		if err != nil {
			return nil, err
		}
		review.Images = string(imagesJSON)
	}

	err = s.reviewRepo.Update(review)
	if err != nil {
		return nil, err
	}

	return review, nil
}

// DeleteReview deletes a review
func (s *ReviewService) DeleteReview(reviewID, touristID uint) error {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	// Verify ownership
	if review.TouristID != touristID {
		return errors.New("unauthorized to delete this review")
	}

	return s.reviewRepo.Delete(reviewID)
}

// GetTourRatingStats returns rating statistics for a tour
func (s *ReviewService) GetTourRatingStats(tourID uint) (map[string]interface{}, error) {
	avgRating, err := s.reviewRepo.GetAverageRating(tourID)
	if err != nil {
		return nil, err
	}

	reviewCount, err := s.reviewRepo.GetReviewCount(tourID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"averageRating": avgRating,
		"reviewCount":   reviewCount,
	}, nil
}

// getUserUsername fetches username from stakeholders service
func (s *ReviewService) getUserUsername(userID uint) (string, error) {
	url := fmt.Sprintf("http://stakeholders-service:8080/api/v1/users/batch?ids=%d", userID)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("user not found: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var users []struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	}

	if err := json.Unmarshal(body, &users); err != nil {
		return "", fmt.Errorf("failed to parse users: %w", err)
	}

	if len(users) == 0 {
		return "", fmt.Errorf("user not found")
	}

	return users[0].Username, nil
}
