package repository

import (
	"time"
	"tour-service/internal/models"

	"gorm.io/gorm"
)

// TourRepository brine o komunikaciji sa bazom
type TourRepository struct {
	DB *gorm.DB
}

// NewTourRepository kreira novu instancu repozitorijuma
func NewTourRepository(db *gorm.DB) *TourRepository {
	return &TourRepository{DB: db}
}

func (r *TourRepository) Create(tour *models.Tour) error {
	return r.DB.Create(tour).Error
}

func (r *TourRepository) FindByAuthorID(authorID uint) ([]models.Tour, error) {
	var tours []models.Tour
	if err := r.DB.Where("author_id = ?", authorID).Find(&tours).Error; err != nil {
		return nil, err
	}
	return tours, nil
}

// FindByID finds tour by ID
func (r *TourRepository) FindByID(tourID uint) (*models.Tour, error) {
	var tour models.Tour
	if err := r.DB.First(&tour, tourID).Error; err != nil {
		return nil, err
	}
	return &tour, nil
}

// FindByIDWithRelations finds tour with keypoints and durations
func (r *TourRepository) FindByIDWithRelations(tourID uint) (*models.Tour, error) {
	var tour models.Tour
	if err := r.DB.Preload("KeyPoints").Preload("Durations").First(&tour, tourID).Error; err != nil {
		return nil, err
	}
	return &tour, nil
}

// CreateDuration creates a new tour duration
func (r *TourRepository) CreateDuration(duration *models.TourDuration) error {
	return r.DB.Create(duration).Error
}

// PublishTour updates tour status to published
func (r *TourRepository) PublishTour(tourID uint) error {
	now := time.Now()
	return r.DB.Model(&models.Tour{}).Where("id = ?", tourID).Updates(map[string]interface{}{
		"status":       models.Published,
		"published_at": now,
	}).Error
}

// ArchiveTour updates tour status to archived
func (r *TourRepository) ArchiveTour(tourID uint) error {
	now := time.Now()
	return r.DB.Model(&models.Tour{}).Where("id = ?", tourID).Updates(map[string]interface{}{
		"status":      models.Archived,
		"archived_at": now,
	}).Error
}

// ActivateTour reactivates archived tour (back to published)
func (r *TourRepository) ActivateTour(tourID uint) error {
	return r.DB.Model(&models.Tour{}).Where("id = ?", tourID).Updates(map[string]interface{}{
		"status":      models.Published,
		"archived_at": nil,
	}).Error
}

// UpdateDistance updates tour distance
func (r *TourRepository) UpdateDistance(tourID uint, distance float64) error {
	return r.DB.Model(&models.Tour{}).Where("id = ?", tourID).Update("distance", distance).Error
}

// FindAllPublished finds all published tours with their first keypoint only
func (r *TourRepository) FindAllPublished() ([]models.Tour, error) {
	var tours []models.Tour
	// Get all published tours with only the first keypoint (order = 1)
	if err := r.DB.Preload("KeyPoints", "\"order\" = 1").
		Where("status = ?", models.Published).
		Find(&tours).Error; err != nil {
		return nil, err
	}
	return tours, nil
}
