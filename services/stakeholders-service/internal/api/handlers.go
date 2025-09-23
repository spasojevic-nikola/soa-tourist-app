package api

import (
	"encoding/json"
	"net/http"
	"stakeholders-service/internal/models"
	"time"

	"gorm.io/gorm"
)

// Handler struktura čuva zavisnosti, kao što je konekcija sa bazom.
type Handler struct {
	DB *gorm.DB
}

// NewHandler kreira novu instancu Handler-a.
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

// GetProfile dobavlja profil ulogovanog korisnika.
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateProfile ažurira profil ulogovanog korisnika.
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Ažuriranje dozvoljenih polja
	if firstName, ok := updateData["first_name"].(string); ok {
		user.FirstName = firstName
	}
	if lastName, ok := updateData["last_name"].(string); ok {
		user.LastName = lastName
	}
	if profileImage, ok := updateData["profile_image"].(string); ok {
		user.ProfileImage = profileImage
	}
	if biography, ok := updateData["biography"].(string); ok {
		user.Biography = biography
	}
	if motto, ok := updateData["motto"].(string); ok {
		user.Motto = motto
	}

	user.UpdatedAt = time.Now()
	if err := h.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser kreira novog korisnika (obično pozvano od strane auth-service).
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetAllUsers (Admin only) - primer kako bi izgledao hendler za admina.
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := h.DB.Find(&users).Error; err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}