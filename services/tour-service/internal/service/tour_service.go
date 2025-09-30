package service

import (
	"errors"
	"tour-service/internal/dto"
	"tour-service/internal/models"
	"tour-service/internal/repository"
)


type TourService struct {
	Repo *repository.TourRepository
}

// NewTourService kreira novu instancu servisa
func NewTourService(repo *repository.TourRepository) *TourService {
	return &TourService{Repo: repo}
}

// CreateTour obrađuje DTO, primenjuje pravila i poziva repozitorijum
// CreateTour kreira turu SA keypointsima
func (s *TourService) CreateTour(authorID uint, req dto.CreateTourRequest) (*models.Tour, error) {
	// Validacija
	if req.Name == "" {
		return nil, errors.New("tour name is required")
	}
	if len(req.KeyPoints) == 0 {
		return nil, errors.New("at least one key point is required")
	}

	// Kreiraj turu
	tour := &models.Tour{
		AuthorID:    authorID,
		Name:        req.Name,
		Description: req.Description,
		Difficulty:  models.TourDifficulty(req.Difficulty),
		Tags:        req.Tags,
		Status:      models.Draft,
		Price:       0,
	}

	err := s.Repo.Create(tour)
	if err != nil {
		return nil, err
	}

	// Kreiraj key points
	keyPointRepo := repository.NewKeyPointRepository(s.Repo.DB)
	for i, kpReq := range req.KeyPoints {
		keyPoint := &models.KeyPoint{
			TourID:      tour.ID,
			Name:        kpReq.Name,
			Description: kpReq.Description,
			Latitude:    kpReq.Latitude,
			Longitude:   kpReq.Longitude,
			Image:       kpReq.Image,
			Order:       i + 1,
		}
		err = keyPointRepo.Create(keyPoint)
		if err != nil {
			return nil, err
		}
	}

	return tour, nil
}


func (s *TourService) GetToursByAuthor(authorID uint) ([]models.Tour, error) {
	return s.Repo.FindByAuthorID(authorID)
}