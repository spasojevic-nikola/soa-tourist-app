package repository

import (
	"tour-service/internal/models"

	"gorm.io/gorm"
)

type KeyPointRepository struct {
	DB *gorm.DB
}

func NewKeyPointRepository(db *gorm.DB) *KeyPointRepository {
	return &KeyPointRepository{DB: db}
}

// Create creates a new key point
func (r *KeyPointRepository) Create(keyPoint *models.KeyPoint) error {
	return r.DB.Create(keyPoint).Error
}

// FindByTourID finds all key points for a specific tour
func (r *KeyPointRepository) FindByTourID(tourID uint) ([]models.KeyPoint, error) {
	var keyPoints []models.KeyPoint
	err := r.DB.Where("tour_id = ?", tourID).Order("\"order\" ASC").Find(&keyPoints).Error
	return keyPoints, err
}

// FindByID finds a key point by ID
func (r *KeyPointRepository) FindByID(id uint) (*models.KeyPoint, error) {
	var keyPoint models.KeyPoint
	err := r.DB.First(&keyPoint, id).Error
	return &keyPoint, err
}

// Update updates a key point
func (r *KeyPointRepository) Update(keyPoint *models.KeyPoint) error {
	return r.DB.Save(keyPoint).Error
}

// Delete deletes a key point
func (r *KeyPointRepository) Delete(id uint) error {
	return r.DB.Delete(&models.KeyPoint{}, id).Error
}

// DeleteByTourID deletes all key points for a specific tour
func (r *KeyPointRepository) DeleteByTourID(tourID uint) error {
	return r.DB.Where("tour_id = ?", tourID).Delete(&models.KeyPoint{}).Error
}