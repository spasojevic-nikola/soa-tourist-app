package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	Role      string    `json:"role" gorm:"not null;default:'tourist'"`
	Blocked   bool      `json:"blocked" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the table name for the User model
func (User) TableName() string {
	return "auth_users"
}

var db *gorm.DB
var jwtKey = []byte("super-secret-key")

// JWT claims
type Claims struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type Handler struct {
	DB *gorm.DB
}

// Init DB
func initDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("Successfully connected to the database!")
			break // Izlazi iz petlje ako je konekcija uspela
		}
		fmt.Printf("Failed to connect to database. Retrying in %d seconds... Error: %v\n", 5, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Fatal: Could not connect to the database after multiple retries. Error: %v\n", err)
	}

	// Auto migrate
	db.Table("auth_users").AutoMigrate(&User{})
}

// Register
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// hesiranje lozinke
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "User already exists or DB error", http.StatusConflict)
		return
	}

	// **Dodatna logika za slanje podataka na stakeholders-service**
	stakeholderUser := StakeholderUser{
		ID:           user.ID,
		FirstName:    "",
		LastName:     "",
		Username:     user.Username,
		Role:         user.Role,
		ProfileImage: "",
		Biography:    "",
		Motto:        "",
	}

	jsonStakeholderUser, err := json.Marshal(stakeholderUser)
	if err != nil {
		log.Println("Failed to marshal stakeholder user:", err)
	}

	// Slati POST zahtjev na stakeholders-service
	resp, err := http.Post("http://stakeholders-service:8080/api/v1/user", "application/json", bytes.NewBuffer(jsonStakeholderUser))
	if err != nil {
		log.Println("Failed to create user in stakeholders-service:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to create user in stakeholders-service. Status: %d", resp.StatusCode)
	}

	// Kreiranje tokena
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"accessToken": tokenString})
}

// Login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var user User
	if err := db.Where("username = ?", creds.Username).First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Poredjenje lozinke iz zahtjeva sa hesiranom lozinkom iz baze
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check if user is blocked
	if user.Blocked {
		http.Error(w, "Account is blocked", http.StatusForbidden)
		return
	}

	// Kreiranje tokena i vracanje odgovora
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"accessToken": tokenString})
}

// Health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "auth-service healthy"})
}

// Admin middleware to check if user is admin
func adminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		if claims.Role != "administrator" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// Get all users (admin only)
func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Remove passwords from response
	for i := range users {
		users[i].Password = ""
	}

	json.NewEncoder(w).Encode(users)
}

// Block user (admin only)
func blockUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.Role == "administrator" {
		http.Error(w, "Cannot block administrator", http.StatusBadRequest)
		return
	}

	user.Blocked = true
	user.UpdatedAt = time.Now()

	if err := db.Save(&user).Error; err != nil {
		http.Error(w, "Failed to block user", http.StatusInternalServerError)
		return
	}

	user.Password = "" // Remove password from response
	json.NewEncoder(w).Encode(user)
}

// Unblock user (admin only)
func unblockUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.Blocked = false
	user.UpdatedAt = time.Now()

	if err := db.Save(&user).Error; err != nil {
		http.Error(w, "Failed to unblock user", http.StatusInternalServerError)
		return
	}

	user.Password = "" // Remove password from response
	json.NewEncoder(w).Encode(user)
}

// GetUserByID - vraća osnovne podatke o korisniku
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := h.DB.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"blocked":  user.Blocked,
	})
}

func main() {
	initDB()

	handler := &Handler{DB: db}

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1/auth").Subrouter()

	// Public routes
	api.HandleFunc("/register", registerHandler).Methods("POST")
	api.HandleFunc("/login", loginHandler).Methods("POST")

	// Admin routes (protected)
	admin := api.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/users", adminMiddleware(getAllUsersHandler)).Methods("GET")
	admin.HandleFunc("/users/{id}/block", adminMiddleware(blockUserHandler)).Methods("POST")
	admin.HandleFunc("/users/{id}/unblock", adminMiddleware(unblockUserHandler)).Methods("POST")

	// ruta za dohvatanje korisnika po ID
	api.HandleFunc("/user/{id}", handler.GetUserByID).Methods("GET") // Koristimo handler.GetUserByID

	r.HandleFunc("/health", healthHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084" //promenjeno!!
	}

	fmt.Println("Auth service running on port", port)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:4200"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS", "DELETE"})

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
