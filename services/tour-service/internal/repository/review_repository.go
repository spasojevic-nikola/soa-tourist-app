package repository

import (
	"tour-service/internal/models"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// Create creates a new review
func (r *ReviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

// FindByID finds a review by ID
func (r *ReviewRepository) FindByID(id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// FindByTourID finds all reviews for a specific tour
func (r *ReviewRepository) FindByTourID(tourID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("tour_id = ?", tourID).Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

// FindByTouristID finds all reviews by a specific tourist
func (r *ReviewRepository) FindByTouristID(touristID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("tourist_id = ?", touristID).Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

// Update updates a review
func (r *ReviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

// Delete deletes a review
func (r *ReviewRepository) Delete(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}

// GetAverageRating calculates the average rating for a tour
func (r *ReviewRepository) GetAverageRating(tourID uint) (float64, error) {
	var avgRating float64
	err := r.db.Model(&models.Review{}).
		Where("tour_id = ?", tourID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avgRating).Error
	return avgRating, err
}

// GetReviewCount returns the count of reviews for a tour
func (r *ReviewRepository) GetReviewCount(tourID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Review{}).Where("tour_id = ?", tourID).Count(&count).Error
	return count, err
}
