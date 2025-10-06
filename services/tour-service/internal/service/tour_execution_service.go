package service

import (
	"errors"
	"log"
	"math"
	"time"
	"tour-service/internal/models"
	"tour-service/internal/repository"
)

type TourExecutionService struct {
	repo *repository.TourExecutionRepository
	purchaseChecker PurchaseChecker
}

func NewTourExecutionService(repo *repository.TourExecutionRepository, purchaseChecker PurchaseChecker) *TourExecutionService {
	return &TourExecutionService{
		repo: repo,
		purchaseChecker: purchaseChecker,
	}
}

func (s *TourExecutionService) StartTour(tourID uint, touristID uint, startLat, startLng float64) (*models.TourExecution, error) {
	// Provera da li već postoji aktivna sesija
	active, _ := s.repo.GetActiveExecution(touristID, tourID)
	if active != nil {
		return nil, errors.New("tour already in progress")
	}

	// PROVERA KUPOVINE - poziv ka purchase servisu
	hasPurchased, err := s.purchaseChecker.HasPurchasedTour(touristID, tourID)
	if err != nil {
		// Ako purchase servis nije dostupan, privremeno dozvoli za testiranje
		// U produkciji ovo treba da vrati grešku
		log.Printf("Purchase service unavailable: %v - temporarily allowing tour start", err)
		// return nil, fmt.Errorf("purchase service unavailable: %v", err)
	} else if !hasPurchased {
		return nil, errors.New("must purchase tour before starting")
	}

	execution := &models.TourExecution{
		TourID:             tourID,
		TouristID:          touristID,
		Status:             models.ExecutionStarted,
		StartTime:          time.Now(),
		LastActivity:       time.Now(),
		CompletedKeyPoints: []uint{},
		StartingLatitude:   startLat,
		StartingLongitude:  startLng,
	}

	return s.repo.Create(execution)
}

func (s *TourExecutionService) CheckPosition(executionID uint, currentLat, currentLng float64) ([]uint, error) {
	// 1. Dohvati execution i turu
	execution, err := s.repo.GetByID(executionID)
	if err != nil {
		return nil, err
	}

	// 2. Dohvati sve key pointove za turu
	keyPoints, err := s.repo.GetKeyPointsByTour(execution.TourID)
	if err != nil {
		return nil, err
	}

	var newlyCompleted []uint

	// 3. Proveri svaki key point
	for _, kp := range keyPoints {
		// Preskoči već završene
		if contains(execution.CompletedKeyPoints, kp.ID) {
			continue
		}

		// 4. Izračunaj distancu do key pointa
		distance := calculateDistance(currentLat, currentLng, kp.Latitude, kp.Longitude)
		
		// 5. Ako je blizu (npr. unutar 50m), označi kao završen
		if distance <= 0.05 { // 50 metara
			execution.CompletedKeyPoints = append(execution.CompletedKeyPoints, kp.ID)
			newlyCompleted = append(newlyCompleted, kp.ID)
		}
	}

	// 6. Ažuriraj execution sa novim completed key points
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

// Pomocne funkcije
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// Haversine formula implementacija
	const R = 6371 // Earth radius in kilometers
	
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c // Distance in kilometers
	return distance
}

func contains(slice []uint, item uint) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}