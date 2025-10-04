package main

import (
	"fmt"
	"log"
	"net/http"

	"shopping-cart-service/internal/api"
	"shopping-cart-service/internal/database"
	"shopping-cart-service/internal/repository"
	"shopping-cart-service/internal/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// 1. Inicijalizacija Baze
	mongoDB := database.InitDB() 

	// 2. Inicijalizacija Središnjih slojeva
	cartRepo := repository.NewCartRepository(mongoDB)
	cartService := service.NewCartService(cartRepo)
	cartHandler := api.NewHandler(cartService) 

	// 3. Postavljanje Routera
	r := mux.NewRouter()

	// CORS opcije
	corsOpts := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}), 
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-User-ID"}),
	)

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

	// 4. Pokretanje servera
	fmt.Println("Shopping Cart service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", corsOpts(r)))
}