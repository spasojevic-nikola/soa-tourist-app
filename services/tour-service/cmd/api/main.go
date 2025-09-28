package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"tour-service/internal/api"
	"tour-service/internal/database"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db := database.InitDB()
	apiHandler := api.NewHandler(db)

	r := mux.NewRouter()
	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	// Ruta za kreiranje ture
	apiV1.Handle("/tours", api.AuthMiddleware(api.AuthorOrAdminAuthMiddleware(apiHandler.CreateTour))).Methods("POST")

	// Health check ruta
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")
	
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(r)

	fmt.Println("Tour service running on internal port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}