package service

import (
	"errors"
	"log"
	"math"
	"time"
	"tour-service/internal/models"
	"tour-service/internal/repository"

	"github.com/lib/pq"
)

type TourExecutionService struct {
	repo            *repository.TourExecutionRepository
	purchaseChecker PurchaseChecker
}

func NewTourExecutionService(repo *repository.TourExecutionRepository, purchaseChecker PurchaseChecker) *TourExecutionService {
	return &TourExecutionService{
		repo:            repo,
		purchaseChecker: purchaseChecker,
	}
}

func (s *TourExecutionService) StartTour(tourID uint, touristID uint, startLat, startLng float64) (*models.TourExecution, error) {
	active, _ := s.repo.GetActiveExecution(touristID, tourID)
	if active != nil {
		return nil, errors.New("tour already in progress")
	}

	hasPurchased, err := s.purchaseChecker.HasPurchasedTour(touristID, tourID)
	if err != nil {
		log.Printf("Purchase service unavailable: %v - temporarily allowing tour start", err)
	} else if !hasPurchased {
		return nil, errors.New("must purchase tour before starting")
	}

	execution := &models.TourExecution{
		TourID:             tourID,
		TouristID:          touristID,
		Status:             models.ExecutionStarted,
		StartTime:          time.Now(),
		LastActivity:       time.Now(),
		CompletedKeyPoints: pq.Int64Array{},
		StartingLatitude:   startLat,
		StartingLongitude:  startLng,
	}

	return s.repo.Create(execution)
}

func (s *TourExecutionService) CheckPosition(executionID uint, currentLat, currentLng float64) ([]int, error) {
	execution, err := s.repo.GetByID(executionID)
	if err != nil {
		return nil, err
	}

	keyPoints, err := s.repo.GetKeyPointsByTour(execution.TourID)
	if err != nil {
		return nil, err
	}

	var newlyCompleted []int

	for _, kp := range keyPoints {
		if contains(execution.CompletedKeyPoints, int64(kp.ID)) {
			continue
		}

		distance := calculateDistance(currentLat, currentLng, kp.Latitude, kp.Longitude)

		if distance <= 0.05 {
			execution.CompletedKeyPoints = append(execution.CompletedKeyPoints, int64(kp.ID))
			newlyCompleted = append(newlyCompleted, int(kp.ID))
		}
	}

	if len(newlyCompleted) > 0 {
		execution.LastActivity = time.Now()
		_, err = s.repo.Update(execution)
		if err != nil {
			return nil, err
		}
	}

	return newlyCompleted, nil
}

func (s *TourExecutionService) CompleteTour(executionID uint) error {
	execution, err := s.repo.GetByID(executionID)
	if err != nil {
		return err
	}

	endTime := time.Now()
	execution.Status = models.ExecutionCompleted
	execution.EndTime = &endTime
	execution.LastActivity = endTime

	_, err = s.repo.Update(execution)
	return err
}

func (s *TourExecutionService) AbandonTour(executionID uint) error {
	execution, err := s.repo.GetByID(executionID)
	if err != nil {
		return err
	}

	endTime := time.Now()
	execution.Status = models.ExecutionAbandoned
	execution.EndTime = &endTime
	execution.LastActivity = endTime

	_, err = s.repo.Update(execution)
	return err
}

func (s *TourExecutionService) GetActiveExecution(touristID uint, tourID uint) (*models.TourExecution, error) {
	return s.repo.GetActiveExecution(touristID, tourID)
}

func (s *TourExecutionService) GetExecutionDetails(executionID uint) (*models.TourExecution, error) {
	return s.repo.GetByID(executionID)
}

func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return distance
}

func contains(slice pq.Int64Array, item int64) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (s *TourExecutionService) GetExecutionsByTour(tourID uint) ([]models.TourExecution, error) {
    return s.repo.GetExecutionsByTour(tourID)
}