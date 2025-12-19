package service

import (
	"errors"
	"math"
	"fmt"
	"tour-service/internal/dto"
	"tour-service/internal/models"
	"tour-service/internal/repository"
	"tour-service/internal/interfaces" 
)

type TourService struct {
	Repo *repository.TourRepository
	PurchaseChecker interfaces.PurchaseChecker
}

// kreira novu instancu servisa
func NewTourService(repo *repository.TourRepository, checker interfaces.PurchaseChecker) *TourService { 
	return &TourService{
	Repo: repo,
	PurchaseChecker: checker,
	}
}

// obrađuje DTO, primenjuje pravila i poziva repozitorijum
// kreira turu SA keypointsima
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

	// racuna distancu ako ima 2+ keypointa
	if len(req.KeyPoints) >= 2 {
		s.CalculateDistance(tour.ID)
	}

	return tour, nil
}

func (s *TourService) GetToursByAuthor(authorID uint) ([]models.Tour, error) {
	return s.Repo.FindByAuthorID(authorID)
}

// vraca sve publishovane ture (vidljive svim korisnicima)
func (s *TourService) GetAllPublishedTours() ([]models.Tour, error) {
	return s.Repo.FindAllPublished()
}

// vraca turu po id sa svim relacijama
func (s *TourService) GetTourByID(tourID, userID uint, authHeader string) (*models.Tour, error) {
	fmt.Println("------------------------------------------")
	fmt.Printf(">>> PROVERA TURE: ID=%d, za KORISNIKA: ID=%d\n", tourID, userID)

	// 1. Uvek prvo dohvati turu sa SVIM ključnim tačkama iz baze
	tour, err := s.Repo.FindByIDWithRelations(tourID)
	if err != nil {
		return nil, err
	}

	if len(tour.KeyPoints) == 0 {
		fmt.Println("!!! Tura nema ključne tačke, vraćam odmah.")
		return tour, nil
	}

	// 2. Proveri da li je korisnik kupio turu (ili je on autor ture)
	isAuthor := tour.AuthorID == userID
	fmt.Printf(">>> DA LI JE AUTOR? %t\n", isAuthor)

	fmt.Println(">>> POZIVAM SHOPPING-CART-SERVICE DA PROVERIM KUPOVINU...")
	hasPurchased, err := s.PurchaseChecker.HasUserPurchasedTour(userID, tourID, authHeader)
	if err != nil {
		fmt.Printf("!!! GREŠKA pri proveri kupovine: %v\n", err)
	}
	
	fmt.Printf(">>> REZULTAT PROVERE KUPOVINE: %t\n", hasPurchased) // <-- KLJUČNI ISPIS

	// 3. Ako korisnik NIJE autor I NIJE kupio turu, onda filtriraj
	if !isAuthor && !hasPurchased {
		fmt.Println(">>> USLOV ISPUNJEN: Korisnik nije autor i nije kupio turu. FILTRIRAM ključne tačke.")
		// Vrati samo prvu ključnu tačku kao "preview"
		tour.KeyPoints = tour.KeyPoints[:1]
	} else {
		fmt.Println(">>> USLOV NIJE ISPUNJEN: Vraćam sve ključne tačke.")
	}
	
	fmt.Println("------------------------------------------")
	// 4. Vrati (potencijalno modifikovanu) turu
	return tour, nil
}

// dodaje trajanje u turu
func (s *TourService) AddDuration(tourID uint, authorID uint, req dto.AddDurationRequest) (*models.TourDuration, error) {
	// verifikuj vlasnistvo ture
	tour, err := s.Repo.FindByID(tourID)
	if err != nil {
		return nil, errors.New("tour not found")
	}
	if tour.AuthorID != authorID {
		return nil, errors.New("unauthorized: not tour author")
	}

	// kreira trajanje
	duration := &models.TourDuration{
		TourID:        tourID,
		TransportType: models.TransportType(req.TransportType),
		DurationMin:   req.DurationMin,
	}

	err = s.Repo.CreateDuration(duration)
	if err != nil {
		return nil, err
	}

	return duration, nil
}

// publishuje turu posle validacije
func (s *TourService) PublishTour(tourID uint, authorID uint) (*models.Tour, error) {
	tour, err := s.Repo.FindByIDWithRelations(tourID)
	if err != nil {
		return nil, errors.New("tour not found")
	}
	if tour.AuthorID != authorID {
		return nil, errors.New("unauthorized: not tour author")
	}
	if tour.Status != models.Draft {
		return nil, errors.New("only draft tours can be published")
	}

	// pravila za validaciju
	if tour.Name == "" || tour.Description == "" {
		return nil, errors.New("tour must have name and description")
	}
	if len(tour.KeyPoints) < 2 {
		return nil, errors.New("tour must have at least 2 key points")
	}

	// publishuje turu
	err = s.Repo.PublishTour(tourID)
	if err != nil {
		return nil, err
	}

	return s.Repo.FindByIDWithRelations(tourID)
}

// arhivira publishovanu turu
func (s *TourService) ArchiveTour(tourID uint, authorID uint) (*models.Tour, error) {
	tour, err := s.Repo.FindByID(tourID)
	if err != nil {
		return nil, errors.New("tour not found")
	}
	if tour.AuthorID != authorID {
		return nil, errors.New("unauthorized: not tour author")
	}
	if tour.Status != models.Published {
		return nil, errors.New("only published tours can be archived")
	}

	err = s.Repo.ArchiveTour(tourID)
	if err != nil {
		return nil, err
	}

	return s.Repo.FindByID(tourID)
}

// reacktivira arhiviranu turu
func (s *TourService) ActivateTour(tourID uint, authorID uint) (*models.Tour, error) {
	tour, err := s.Repo.FindByID(tourID)
	if err != nil {
		return nil, errors.New("tour not found")
	}
	if tour.AuthorID != authorID {
		return nil, errors.New("unauthorized: not tour author")
	}
	if tour.Status != models.Archived {
		return nil, errors.New("only archived tours can be activated")
	}

	err = s.Repo.ActivateTour(tourID)
	if err != nil {
		return nil, err
	}

	return s.Repo.FindByID(tourID)
}

// racuna distancu po keypointima
func (s *TourService) CalculateDistance(tourID uint) error {
	tour, err := s.Repo.FindByIDWithRelations(tourID)
	if err != nil {
		return err
	}

	if len(tour.KeyPoints) < 2 {
		return nil //nema distance da se racuna
	}

	totalDistance := 0.0
	for i := 0; i < len(tour.KeyPoints)-1; i++ {
		kp1 := tour.KeyPoints[i]
		kp2 := tour.KeyPoints[i+1]
		distance := haversineDistance(kp1.Latitude, kp1.Longitude, kp2.Latitude, kp2.Longitude)
		totalDistance += distance
	}

	return s.Repo.UpdateDistance(tourID, totalDistance)
}

// haversineDistance calculates distance between two lat/long points in kilometers
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0 // Earth radius u kilometrima

	// konvertuje stepene u radijane
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	// Haversine formula
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Asin(math.Sqrt(a))

	return earthRadius * c
}
