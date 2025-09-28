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
func (s *TourService) CreateTour(authorID uint, req dto.CreateTourRequest) (*models.Tour, error) {
	// Validacija (ovo se može proširiti)
	if req.Name == "" {
		return nil, errors.New("tour name is required")
	}

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

	return tour, nil
}

func (s *TourService) GetToursByAuthor(authorID uint) ([]models.Tour, error) {
	return s.Repo.FindByAuthorID(authorID)
}