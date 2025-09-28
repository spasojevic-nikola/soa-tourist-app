package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"tour-service/internal/api"
	"tour-service/internal/database"
	"tour-service/internal/repository" 
	"tour-service/internal/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db := database.InitDB()
	
	tourRepo := repository.NewTourRepository(db)
	tourService := service.NewTourService(tourRepo)
	apiHandler := api.NewHandler(tourService)

	r := mux.NewRouter()
	apiV1 := r.PathPrefix("/api/v1/tours").Subrouter()

	
	apiV1.Handle("/create-tour", api.AuthMiddleware(api.AuthorOrAdminAuthMiddleware(apiHandler.CreateTour))).Methods("POST")
	apiV1.Handle("", api.AuthMiddleware(api.AuthorOrAdminAuthMiddleware(apiHandler.GetMyTours))).Methods("GET")


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