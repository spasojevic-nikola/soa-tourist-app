package api

import (
	"encoding/json"
	"net/http"
	"tour-service/internal/models"
	"tour-service/internal/dto"   
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) CreateTour(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)

	// 2. Dekodiraj telo zahteva u DT
	var req dto.CreateTourRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//  Validacija za Difficulty "enum"
    switch models.TourDifficulty(req.Difficulty) {
    case models.Easy, models.Medium, models.Hard, models.Expert:
    default:
        http.Error(w, "Invalid tour difficulty value. Must be Easy, Medium, Hard, or Expert.", http.StatusBadRequest)
        return
    }

	// 3. Kreiranje  na osnovu DTO 
	tour := models.Tour{
		AuthorID:    userID,
		Name:        req.Name,
		Description: req.Description,
   	 	Difficulty:  models.TourDifficulty(req.Difficulty),
		Tags:        req.Tags,
		Status:      models.Draft,
		Price:       0,            
	}

	// 4. Sacuvaj u bazi
	if err := h.DB.Create(&tour).Error; err != nil {
		http.Error(w, "Failed to create tour", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tour)
}
func (h *Handler) GetMyTours(w http.ResponseWriter, r *http.Request) {
    userID, _ := r.Context().Value("userID").(uint)

    var tours []models.Tour

    if err := h.DB.Where("author_id = ?", userID).Find(&tours).Error; err != nil {
        http.Error(w, "Failed to retrieve tours", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tours)
}