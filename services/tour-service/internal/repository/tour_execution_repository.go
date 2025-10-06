package repository

import (
	"tour-service/internal/models"
	"gorm.io/gorm"
)

type TourExecutionRepository struct {
	DB *gorm.DB
}

func NewTourExecutionRepository(db *gorm.DB) *TourExecutionRepository {
	return &TourExecutionRepository{DB: db}
}

func (r *TourExecutionRepository) Create(execution *models.TourExecution) (*models.TourExecution, error) {
	err := r.DB.Create(execution).Error
	return execution, err
}

func (r *TourExecutionRepository) GetByID(id uint) (*models.TourExecution, error) {
	var execution models.TourExecution
	err := r.DB.First(&execution, id).Error
	return &execution, err
}

func (r *TourExecutionRepository) Update(execution *models.TourExecution) (*models.TourExecution, error) {
	err := r.DB.Save(execution).Error
	return execution, err
}

func (r *TourExecutionRepository) UpdateStatus(id uint, status models.TourExecutionStatus) error {
	return r.DB.Model(&models.TourExecution{}).Where("id = ?", id).Update("status", status).Error
}

func (r *TourExecutionRepository) GetActiveExecution(touristID uint, tourID uint) (*models.TourExecution, error) {
	var execution models.TourExecution
	err := r.DB.Where("tourist_id = ? AND tour_id = ? AND status = ?", touristID, tourID, models.ExecutionStarted).First(&execution).Error
	return &execution, err
}

func (r *TourExecutionRepository) GetKeyPointsByTour(tourID uint) ([]models.KeyPoint, error) {
	var keyPoints []models.KeyPoint
	err := r.DB.Where("tour_id = ?", tourID).Order("\"order\" ASC").Find(&keyPoints).Error
	return keyPoints, err
}