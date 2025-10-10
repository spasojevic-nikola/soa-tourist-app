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
	"tour-service/internal/clients"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db := database.InitDB()

	tourRepo := repository.NewTourRepository(db)
	keyPointRepo := repository.NewKeyPointRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	tourExecutionRepo := repository.NewTourExecutionRepository(db)

	shoppingCartClient := clients.NewShoppingCartClient("http://shopping-cart-service:8081")
	tourService := service.NewTourService(tourRepo, shoppingCartClient)
	keyPointService := service.NewKeyPointService(keyPointRepo, tourRepo)
	reviewService := service.NewReviewService(reviewRepo, tourRepo)
	purchaseChecker := clients.NewRESTPurchaseChecker("http://purchase-service:8082")
	tourExecutionService := service.NewTourExecutionService(tourExecutionRepo, purchaseChecker)

	apiHandler := api.NewHandler(tourService, keyPointService)
	reviewHandler := api.NewReviewHandler(reviewService)
	tourExecutionHandler := api.NewTourExecutionHandler(tourExecutionService)

	r := mux.NewRouter()
	apiV1 := r.PathPrefix("/api/v1/tours").Subrouter()

	// Tour routes
	apiV1.Handle("/create-tour", api.AuthMiddleware(api.AuthorOrAdminAuthMiddleware(apiHandler.CreateTour))).Methods("POST")
	apiV1.Handle("", api.AuthMiddleware(api.AuthorOrAdminAuthMiddleware(apiHandler.GetMyTours))).Methods("GET")
	apiV1.Handle("/published", api.AuthMiddleware(http.HandlerFunc(apiHandler.GetAllPublishedTours))).Methods("GET")
	apiV1.Handle("/{tourId}", api.AuthMiddleware(http.HandlerFunc(apiHandler.GetTourByID))).Methods("GET")
	apiV1.Handle("/{tourId}/publish", api.AuthMiddleware(api.AuthorOrAdminAuthMiddleware(apiHandler.PublishTour))).Methods("PUT")
	apiV1.Handle("/{tourId}/archive", api.AuthMiddleware(api.AuthorOrAdminAuthMiddleware(apiHandler.ArchiveTour))).Methods("PUT")
	apiV1.Handle("/{tourId}/activate", api.AuthMiddleware(api.AuthorOrAdminAuthMiddleware(apiHandler.ActivateTour))).Methods("PUT")
	apiV1.Handle("/{tourId}/duration", api.AuthMiddleware(api.AuthorOrAdminAuthMiddleware(apiHandler.AddDuration))).Methods("POST")

	// KeyPoint routes
	apiV1.Handle("/{tourId}/keypoints", api.AuthMiddleware(api.AuthorOrAdminAuthMiddleware(apiHandler.GetKeyPointsByTour))).Methods("GET")
	apiV1.Handle("/keypoints/{keyPointId}", api.AuthMiddleware(api.AuthorOrAdminAuthMiddleware(apiHandler.UpdateKeyPoint))).Methods("PUT")
	apiV1.Handle("/keypoints/{keyPointId}", api.AuthMiddleware(api.AuthorOrAdminAuthMiddleware(apiHandler.DeleteKeyPoint))).Methods("DELETE")

	// Review routes
	apiV1.Handle("/{tourId}/reviews", api.AuthMiddleware(http.HandlerFunc(reviewHandler.CreateReview))).Methods("POST")
	apiV1.Handle("/{tourId}/reviews", http.HandlerFunc(reviewHandler.GetReviewsByTour)).Methods("GET")
	apiV1.Handle("/{tourId}/reviews/stats", http.HandlerFunc(reviewHandler.GetTourRatingStats)).Methods("GET")
	apiV1.Handle("/reviews/{reviewId}", api.AuthMiddleware(http.HandlerFunc(reviewHandler.UpdateReview))).Methods("PUT")
	apiV1.Handle("/reviews/{reviewId}", api.AuthMiddleware(http.HandlerFunc(reviewHandler.DeleteReview))).Methods("DELETE")
	apiV1.Handle("/my-reviews", api.AuthMiddleware(http.HandlerFunc(reviewHandler.GetMyReviews))).Methods("GET")

	// TourExecution routes
	apiV1.Handle("/{tourId}/start", api.AuthMiddleware(tourExecutionHandler.StartTour)).Methods("POST")
	apiV1.Handle("/executions/{executionId}/check-position", api.AuthMiddleware(tourExecutionHandler.CheckPosition)).Methods("POST")
	apiV1.Handle("/executions/{executionId}/complete", api.AuthMiddleware(tourExecutionHandler.CompleteTour)).Methods("PUT")
	apiV1.Handle("/executions/{executionId}/abandon", api.AuthMiddleware(tourExecutionHandler.AbandonTour)).Methods("PUT")
	apiV1.Handle("/executions/active/{tourId}", api.AuthMiddleware(tourExecutionHandler.GetActiveExecution)).Methods("GET")
	apiV1.Handle("/executions/{executionId}", api.AuthMiddleware(tourExecutionHandler.GetExecutionDetails)).Methods("GET")
	apiV1.Handle("/executions/tour/{tourId}", api.AuthMiddleware(tourExecutionHandler.GetExecutionsByTour)).Methods("GET")

	// Health check ruta
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-User-ID"}),
	)(r)

	fmt.Println("Tour service running on internal port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
