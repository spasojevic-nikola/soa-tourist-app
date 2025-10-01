package main

import (
	"fmt"
	"log"
	"net/http"

	"soa-tourist-app/follower-service/internal/api"
	"soa-tourist-app/follower-service/internal/database"
	"soa-tourist-app/follower-service/internal/repository"
	"soa-tourist-app/follower-service/internal/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Inicijalizacija
	driver := database.InitDB()
	// defer driver.Close(context.Background()) // defer se sada ne koristi jer server radi non-stop

	repo := repository.NewFollowerRepository(driver)
	followerService := service.NewFollowerService(repo)
	handler := api.NewHandler(followerService)

	// Ruter
	r := mux.NewRouter()
	
	// Kreiramo pod-ruter za rute koje zahtevaju autorizaciju
	protectedRoutes := r.PathPrefix("/api/followers").Subrouter()
	protectedRoutes.Use(api.AuthMiddleware) // Primenjujemo middleware na sve rute u ovom pod-ruteru

	// Definišemo rute na pod-ruteru
	protectedRoutes.HandleFunc("/follow/{id:[0-9]+}", handler.Follow).Methods("POST")
	protectedRoutes.HandleFunc("/unfollow/{id:[0-9]+}", handler.Unfollow).Methods("DELETE")
	protectedRoutes.HandleFunc("/check-follow/{id:[0-9]+}", handler.CheckFollow).Methods("GET")
	protectedRoutes.HandleFunc("/recommendations", handler.GetRecommendations).Methods("GET")
	protectedRoutes.HandleFunc("/following", handler.GetFollowingIDs).Methods("GET")
	
	// Podesavanje CORS-a
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}), // Promenite ako je potrebno
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-User-ID"}),
		)(r)

	// Pokretanje servera
	fmt.Println("Follower service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}