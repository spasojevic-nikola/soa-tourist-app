package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"stakeholders-service/internal/dto"
	"stakeholders-service/internal/models"
	"time"

	"gorm.io/gorm"
)

// Handler struktura čuva zavisnosti, kao što je konekcija sa bazom.
type Handler struct {
	DB         *gorm.DB
	authClient *http.Client
}

// NewHandler kreira novu instancu Handler-a.
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		DB:         db,
		authClient: &http.Client{Timeout: 10 * time.Second},
	}
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

// UpdateProfile azurira profil ulogovanog korisnika.
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// 1.Dohvati ID korisnika iz JWT tokena (koji je middleware postavio u kontekst)
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}
	// 2. Pronadji trenutnog korisnika u bazi da bismo imali pocetne podatke
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// 3. Dekodiraj JSON telo zahteva u privremenu mapu, da bismo mogli videti koja polja je korisnik poslao
	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

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
	// 5. Postavi novo vreme azuriranja i sacuvaj u bazi
	user.UpdatedAt = time.Now()
	if err := h.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}
	// 6. Vrati azurirani profil kao odgovor
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
	currentUserID, ok := r.Context().Value("currentUserID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var users []models.User
	if err := h.DB.Find(&users).Error; err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	var userDTOs []dto.UserOverviewDto
	for _, user := range users {
		// preskociti ulogovanog administratora:
		if user.ID == currentUserID {
			continue
		}

		log.Printf("Current User ID from context: %d", currentUserID)

		// pozivanje auth-servisa
		authServiceURL := fmt.Sprintf("http://auth-service:8084/api/v1/auth/user/%d", user.ID)

		resp, err := h.authClient.Get(authServiceURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("Failed to get auth data for user ID %d: %v", user.ID, err)
			continue
		}
		defer resp.Body.Close()

		var authUser struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
			Blocked  bool   `json:"blocked"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&authUser); err != nil {
			log.Printf("Failed to decode auth data for user ID %d: %v", user.ID, err)
			continue
		}

		userDTOs = append(userDTOs, dto.UserOverviewDto{
			ID:           user.ID,
			Username:     authUser.Username,
			Email:        authUser.Email,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			ProfileImage: user.ProfileImage,
			Role:         authUser.Role,
			Blocked:      authUser.Blocked,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userDTOs)
}
