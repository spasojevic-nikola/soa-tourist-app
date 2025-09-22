package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
    "github.com/gorilla/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User model
type User struct {
    ID              uint       `json:"id" gorm:"primaryKey"`
    FirstName       string     `json:"first_name"`
    LastName        string     `json:"last_name"`
    ProfileImage    string     `json:"profile_image"`
    Biography       string     `json:"biography"`
    Motto           string     `json:"motto"`
    Role            string     `json:"role" gorm:"default:'tourist'"`
    IsBlocked       bool       `json:"is_blocked" gorm:"default:false"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`
}

// Database connection
var db *gorm.DB
var jwtKey = []byte("super-secret-key")

// JWT claims
type Claims struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func initDB() {
    dsn := "host=postgres user=postgres password=password dbname=stakeholders port=5432 sslmode=disable"
    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto migrate the schema
    db.Table("stakeholders_users").AutoMigrate(&User{})
}

func main() {
    initDB()

    r := mux.NewRouter()

    // API routes
    api := r.PathPrefix("/api/v1").Subrouter()

    // User registration
    // api.HandleFunc("/register", registerUser).Methods("POST")

    // Admin routes
    //api.HandleFunc("/admin/users", getAllUsers).Methods("GET")
    //api.HandleFunc("/admin/users/{id}/block", blockUser).Methods("POST")
    //api.HandleFunc("/admin/users/{id}/unblock", unblockUser).Methods("POST")

    api.HandleFunc("/user", createUser).Methods("POST")
    
    // Profile routes, now protected by middleware
    api.Handle("/profile", authMiddleware(getProfile)).Methods("GET")
	api.Handle("/profile", authMiddleware(updateProfile)).Methods("PUT")

    // Health check
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
    }).Methods("GET")

    // CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(r)

	fmt.Println("Stakeholders service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

func (User) TableName() string {
    return "stakeholders_users"
}

// Authorization middleware to extract and validate JWT
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userRole", claims.Role)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

// Get all users (Admin only)
func getAllUsers(w http.ResponseWriter, r *http.Request) {
    var users []User
    if err := db.Find(&users).Error; err != nil {
        http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

/*
// Block user (Admin only)
func blockUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    if err := db.Model(&User{}).Where("id = ?", userID).Update("is_blocked", true).Error; err != nil {
        http.Error(w, "Failed to block user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "User blocked successfully"})
}

// Unblock user (Admin only)
func unblockUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    if err := db.Model(&User{}).Where("id = ?", userID).Update("is_blocked", false).Error; err != nil {
        http.Error(w, "Failed to unblock user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "User unblocked successfully"})
}*/

// Get user profile
func getProfile(w http.ResponseWriter, r *http.Request) {
	// Dohvatanje userID-a iz konteksta zahtjeva
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Update user profile
func updateProfile(w http.ResponseWriter, r *http.Request) {
    // Dohvatanje userID-a iz konteksta zahtjeva
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}
    var user User
    if err := db.First(&user, userID).Error; err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    var updateData map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    // Update allowed fields
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
    if err := db.Save(&user).Error; err != nil {
        http.Error(w, "Failed to update profile", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Gorm ce automatski koristiti ID poslat od auth-service
    if err := db.Create(&user).Error; err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}