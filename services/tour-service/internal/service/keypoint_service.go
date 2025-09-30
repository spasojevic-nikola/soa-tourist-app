package service

import (
	"errors"
	"tour-service/internal/dto"
	"tour-service/internal/models"
	"tour-service/internal/repository"
)

type KeyPointService struct {
	KeyPointRepo *repository.KeyPointRepository
	TourRepo     *repository.TourRepository
}

func NewKeyPointService(keyPointRepo *repository.KeyPointRepository, tourRepo *repository.TourRepository) *KeyPointService {
	return &KeyPointService{
		KeyPointRepo: keyPointRepo,
		TourRepo:     tourRepo,
	}
}

// GetKeyPointsByTour gets all key points for a specific tour
func (s *KeyPointService) GetKeyPointsByTour(tourID uint) ([]models.KeyPoint, error) {
	return s.KeyPointRepo.FindByTourID(tourID)
}

// UpdateKeyPoint updates an existing key point
func (s *KeyPointService) UpdateKeyPoint(keyPointID uint, authorID uint, req dto.UpdateKeyPointRequest) (*models.KeyPoint, error) {
	// Get the key point
	keyPoint, err := s.KeyPointRepo.FindByID(keyPointID)
	if err != nil {
		return nil, errors.New("key point not found")
	}

	// Verify that the author owns the tour
	tours, err := s.TourRepo.FindByAuthorID(authorID)
	if err != nil {
		return nil, err
	}

	var hasPermission bool
	for _, tour := range tours {
		if tour.ID == keyPoint.TourID {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return nil, errors.New("you don't have permission to modify this key point")
	}

	// Update fields if provided
	if req.Name != "" {
		keyPoint.Name = req.Name
	}
	if req.Description != "" {
		keyPoint.Description = req.Description
	}
	if req.Latitude != 0 {
		keyPoint.Latitude = req.Latitude
	}
	if req.Longitude != 0 {
		keyPoint.Longitude = req.Longitude
	}
	if req.Image != "" {
		keyPoint.Image = req.Image
	}
	if req.Order != 0 {
		keyPoint.Order = req.Order
	}

	err = s.KeyPointRepo.Update(keyPoint)
	if err != nil {
		return nil, err
	}

	return keyPoint, nil
}

// DeleteKeyPoint deletes a key point
func (s *KeyPointService) DeleteKeyPoint(keyPointID uint, authorID uint) error {
	// Get the key point to check ownership
	keyPoint, err := s.KeyPointRepo.FindByID(keyPointID)
	if err != nil {
		return errors.New("key point not found")
	}

	// Verify that the author owns the tour
	tours, err := s.TourRepo.FindByAuthorID(authorID)
	if err != nil {
		return err
	}

	var hasPermission bool
	for _, tour := range tours {
		if tour.ID == keyPoint.TourID {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return errors.New("you don't have permission to delete this key point")
	}

	return s.KeyPointRepo.Delete(keyPointID)
}