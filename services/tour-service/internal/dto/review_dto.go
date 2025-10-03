package dto

import "time"

// CreateReviewRequest represents the payload for creating a review
type CreateReviewRequest struct {
	TourID    uint      `json:"tourId" binding:"required"`
	Rating    int       `json:"rating" binding:"required,min=1,max=5"`
	Comment   string    `json:"comment"`
	VisitDate time.Time `json:"visitDate" binding:"required"`
	Images    []string  `json:"images"` // Array of image URLs
}

// UpdateReviewRequest represents the payload for updating a review
type UpdateReviewRequest struct {
	Rating    int       `json:"rating" binding:"min=1,max=5"`
	Comment   string    `json:"comment"`
	VisitDate time.Time `json:"visitDate"`
	Images    []string  `json:"images"`
}

// ReviewResponse represents the review response with tourist info
type ReviewResponse struct {
	ID              uint      `json:"id"`
	TourID          uint      `json:"tourId"`
	TouristID       uint      `json:"touristId"`
	TouristName     string    `json:"touristName"`     // From stakeholder service
	TouristUsername string    `json:"touristUsername"` // From stakeholder service
	Rating          int       `json:"rating"`
	Comment         string    `json:"comment"`
	VisitDate       time.Time `json:"visitDate"`
	Images          []string  `json:"images"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
