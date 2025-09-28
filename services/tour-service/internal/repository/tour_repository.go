package repository

import (
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