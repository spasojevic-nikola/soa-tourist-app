package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// Uvozimo naše nove pakete
	"stakeholders-service/internal/api"
	"stakeholders-service/internal/database"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// 1. Inicijalizacija baze pozivanjem funkcije iz `database` paketa
	db := database.InitDB()

	// 2. Kreiranje instance našeg API hendlera i prosleđivanje konekcije
	apiHandler := api.NewHandler(db)

	// 3. Podešavanje rutera
	r := mux.NewRouter()
	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	// 4. Definisanje ruta i povezivanje sa metodama iz apiHandler-a
	apiV1.HandleFunc("/user", apiHandler.CreateUser).Methods("POST")
	
	// Zaštićene rute koriste AuthMiddleware iz `api` paketa
	apiV1.Handle("/profile", api.AuthMiddleware(apiHandler.GetProfile)).Methods("GET")
	apiV1.Handle("/profile", api.AuthMiddleware(apiHandler.UpdateProfile)).Methods("PUT")
	
	// Admin ruta za pregled svih korisnika
	apiV1.Handle("/admin/users", api.AuthMiddleware(api.AdminAuthMiddleware(apiHandler.GetAllUsers))).Methods("GET")

	// Health check ruta može ostati ovde
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	// 5. Podešavanje CORS-a
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(r)

	// 6. Pokretanje servera
	fmt.Println("Stakeholders service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}