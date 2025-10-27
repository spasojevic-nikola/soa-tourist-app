package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"shopping-cart-service/internal/api"
	"shopping-cart-service/internal/database"
	"shopping-cart-service/internal/grpc"
	"shopping-cart-service/internal/repository"
	"shopping-cart-service/internal/service"
	"shopping-cart-service/internal/client"

	"github.com/gorilla/mux"
)

func main() {
	// 1. Inicijalizacija Baze
	mongoDB := database.InitDB() 

	// 2. Inicijalizacija Klijenata
    tourServiceURL := os.Getenv("TOUR_SERVICE_URL")
    if tourServiceURL == "" {
        tourServiceURL = "http://tour-service:8080" // Port 8080 je u tour-service
    }
    tourClient := client.NewTourServiceClient(tourServiceURL)

	// 2. Inicijalizacija Središnjih slojeva
	cartRepo := repository.NewCartRepository(mongoDB)
    cartService := service.NewCartService(cartRepo, tourClient)
	cartHandler := api.NewHandler(cartService) 

	// 3. POKRENI gRPC SERVER U POZADINI 
	go func() {
		log.Println("Starting gRPC server on port 50051...")
		grpc.StartGRPCServer(cartService, "50051")
	}()

	// 4. Postavljanje Routera
	r := mux.NewRouter()

	// CORS removed - requests now go through API gateway which handles CORS

	apiV1 := r.PathPrefix("/api/v1/cart").Subrouter()

	// Rute za korpu (zahtevaju autentikaciju)
	apiV1.HandleFunc("/items/{tourId}", api.AuthMiddleware(cartHandler.RemoveItem)).Methods("DELETE") 
	apiV1.HandleFunc("", api.AuthMiddleware(cartHandler.GetCart)).Methods("GET")
	apiV1.HandleFunc("/items", api.AuthMiddleware(cartHandler.AddItemToCart)).Methods("POST") 
	apiV1.HandleFunc("/checkout", api.AuthMiddleware(cartHandler.Checkout)).Methods("POST")
	//da li korisnik ima token
	apiV1.HandleFunc("/purchase-status/{tourId}", api.AuthMiddleware(cartHandler.HasPurchaseToken)).Methods("GET") 


	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "healthy"}`)
	}).Methods("GET")

	// 5. Pokretanje servera
	fmt.Println("Shopping Cart service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}