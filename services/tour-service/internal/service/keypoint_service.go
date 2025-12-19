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

// GetKeyPointsByTour vraca sve keypointe za turu
func (s *KeyPointService) GetKeyPointsByTour(tourID uint) ([]models.KeyPoint, error) {
	return s.KeyPointRepo.FindByTourID(tourID)
}

// UpdateKeyPoint updatuje postojeci keypoint
func (s *KeyPointService) UpdateKeyPoint(keyPointID uint, authorID uint, req dto.UpdateKeyPointRequest) (*models.KeyPoint, error) {
	// Get the key point
	keyPoint, err := s.KeyPointRepo.FindByID(keyPointID)
	if err != nil {
		return nil, errors.New("key point not found")
	}

	// Verifikuj da autor poseduje turu
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

	// Update polja
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

	// Recalkulise distancu posle azuriranja keypointa (ako se koordinate menjaju)
	if req.Latitude != 0 || req.Longitude != 0 {
		s.calculateAndUpdateDistance(keyPoint.TourID)
	}

	return keyPoint, nil
}

// Obrisi keypoint
func (s *KeyPointService) DeleteKeyPoint(keyPointID uint, authorID uint) error {
	// Get key point da proveri vlasnistvo
	keyPoint, err := s.KeyPointRepo.FindByID(keyPointID)
	if err != nil {
		return errors.New("key point not found")
	}

	// Verifikuj da autor poseduje turu
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

	tourID := keyPoint.TourID
	err = s.KeyPointRepo.Delete(keyPointID)
	if err != nil {
		return err
	}

	// Rekalkulisi distancu posle brisanja keypointa
	s.calculateAndUpdateDistance(tourID)

	return nil
}

// kalkulise i azurira distancu ture
func (s *KeyPointService) calculateAndUpdateDistance(tourID uint) error {
	// vrati sve keypointe za turu
	keyPoints, err := s.KeyPointRepo.FindByTourID(tourID)
	if err != nil {
		return err
	}

	// treba mu min 2 keypointa da bi se izracunala distanca
	if len(keyPoints) < 2 {
		// postavi distancu na 0 ako ima manje od 2 keypointa
		return s.TourRepo.UpdateDistance(tourID, 0)
	}

	// racuna ukupnu distancu izmedju keypointa
	totalDistance := 0.0
	for i := 0; i < len(keyPoints)-1; i++ {
		kp1 := keyPoints[i]
		kp2 := keyPoints[i+1]
		distance := haversineDistance(kp1.Latitude, kp1.Longitude, kp2.Latitude, kp2.Longitude)
		totalDistance += distance
	}

	// azurira distancu ture
	return s.TourRepo.UpdateDistance(tourID, totalDistance)
}
