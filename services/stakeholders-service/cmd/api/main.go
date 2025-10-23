package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"stakeholders-service/internal/api"
	"stakeholders-service/internal/database"

	"github.com/gorilla/mux"
)

func main() {
	db := database.InitDB()
	apiHandler := api.NewHandler(db)

	r := mux.NewRouter()

    // 1. Javna ruta za CreateUser ostaje puna, jer se verovatno ne rutira kroz JWT middleware
    // Ako se ova ruta rutira kroz Gateway: /api/v1/user
	r.HandleFunc("/api/v1/user", apiHandler.CreateUser).Methods("POST")
	
	// Interna komunikacija (bez /api/v1 prefiksa)
	r.HandleFunc("/users/batch", apiHandler.GetUsersBatch).Methods("GET") 

	// KREIRANJE ZAŠTIĆENOG RUTERA
	// Ostavljamo SVE ostale rute bez prefixa /api/v1
	// Koristimo .NewRoute().Subrouter() da ne bismo imali PathPrefix("/api/v1")
	protectedRoutes := r.NewRoute().Subrouter()
	
	// FIX GRESKE KOMPAJLIRANJA (MUX.Use problem)
    // Koristimo klasični AuthMiddleware za sve rute u ovom subrouteru
	protectedRoutes.Use(func(next http.Handler) http.Handler {
		return api.AuthMiddleware(next.ServeHTTP)
	})

	// 4. Definisanje ruta - OVE RUTE SADA OČEKUJU PUTANJE KOJE GATEWAY PROSLEĐUJE (npr. SAMO /profile)
	
	// Rute za profil
	protectedRoutes.HandleFunc("/profile", apiHandler.GetProfile).Methods("GET")
	protectedRoutes.HandleFunc("/profile", apiHandler.UpdateProfile).Methods("PUT")
	
	// Admin ruta (AdminAuthMiddleware je primenjen unutar rutiranja, ne na nivou subrutera)
	protectedRoutes.Handle("/admin/users", api.AdminAuthMiddleware(apiHandler.GetAllUsers)).Methods("GET")

	// Pretraga i GetById
	protectedRoutes.HandleFunc("/users/search", apiHandler.SearchUsers).Methods("GET")
	protectedRoutes.HandleFunc("/users/{id:[0-9]+}", apiHandler.GetUserById).Methods("GET")
    
	// Health check ruta ostaje na glavnom ruteru
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	fmt.Println("Stakeholders service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}