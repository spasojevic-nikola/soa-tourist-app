package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "tour-service/gen/pb-go/tour"
	"tour-service/internal/api"
	"tour-service/internal/clients"
	"tour-service/internal/database"
	tourgrpc "tour-service/internal/grpc"
	"tour-service/internal/repository"
	"tour-service/internal/service"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
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
	purchaseChecker, err := clients.NewGRPCPurchaseChecker("shopping-cart-service:50051")
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	tourExecutionService := service.NewTourExecutionService(tourExecutionRepo, purchaseChecker)

	apiHandler := api.NewHandler(tourService, keyPointService)
	reviewHandler := api.NewReviewHandler(reviewService)
	tourExecutionHandler := api.NewTourExecutionHandler(tourExecutionService)

	r := mux.NewRouter()
	apiV1 := r.PathPrefix("/api/v1/tours").Subrouter()

	// Tour routes - API Gateway sada radi JWT validaciju i prosleđuje X-User-* headere
	apiV1.HandleFunc("/create-tour", apiHandler.CreateTour).Methods("POST")
	apiV1.HandleFunc("", apiHandler.GetMyTours).Methods("GET")
	apiV1.HandleFunc("/published", apiHandler.GetAllPublishedTours).Methods("GET")
	apiV1.HandleFunc("/{tourId}", apiHandler.GetTourByID).Methods("GET")
	apiV1.HandleFunc("/{tourId}/publish", apiHandler.PublishTour).Methods("PUT")
	apiV1.HandleFunc("/{tourId}/archive", apiHandler.ArchiveTour).Methods("PUT")
	apiV1.HandleFunc("/{tourId}/activate", apiHandler.ActivateTour).Methods("PUT")
	apiV1.HandleFunc("/{tourId}/duration", apiHandler.AddDuration).Methods("POST")

	// KeyPoint routes
	apiV1.HandleFunc("/{tourId}/keypoints", apiHandler.GetKeyPointsByTour).Methods("GET")
	apiV1.HandleFunc("/keypoints/{keyPointId}", apiHandler.UpdateKeyPoint).Methods("PUT")
	apiV1.HandleFunc("/keypoints/{keyPointId}", apiHandler.DeleteKeyPoint).Methods("DELETE")

	// Review routes
	apiV1.HandleFunc("/{tourId}/reviews", reviewHandler.CreateReview).Methods("POST")
	apiV1.HandleFunc("/{tourId}/reviews", reviewHandler.GetReviewsByTour).Methods("GET")
	apiV1.HandleFunc("/{tourId}/reviews/stats", reviewHandler.GetTourRatingStats).Methods("GET")
	apiV1.HandleFunc("/reviews/{reviewId}", reviewHandler.UpdateReview).Methods("PUT")
	apiV1.HandleFunc("/reviews/{reviewId}", reviewHandler.DeleteReview).Methods("DELETE")
	apiV1.HandleFunc("/my-reviews", reviewHandler.GetMyReviews).Methods("GET")

	// TourExecution routes
	apiV1.HandleFunc("/{tourId}/start", tourExecutionHandler.StartTour).Methods("POST")
	apiV1.HandleFunc("/executions/{executionId}/check-position", tourExecutionHandler.CheckPosition).Methods("POST")
	apiV1.HandleFunc("/executions/{executionId}/complete", tourExecutionHandler.CompleteTour).Methods("PUT")
	apiV1.HandleFunc("/executions/{executionId}/abandon", tourExecutionHandler.AbandonTour).Methods("PUT")
	apiV1.HandleFunc("/executions/active/{tourId}", tourExecutionHandler.GetActiveExecution).Methods("GET")
	apiV1.HandleFunc("/executions/{executionId}", tourExecutionHandler.GetExecutionDetails).Methods("GET")
	apiV1.HandleFunc("/executions/tour/{tourId}", tourExecutionHandler.GetExecutionsByTour).Methods("GET")

	// Health check ruta
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	//corsHandler := handlers.CORS(
	//	handlers.AllowedOrigins([]string{"http://localhost:4200"}),
	//	handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
	//	handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-User-ID", "X-User-Username", "X-User-Role"}),
	//)(r)

	// Pokreni gRPC server u goroutine
	go func() {
		grpcServer := grpc.NewServer()
		tourGRPCServer := tourgrpc.NewTourGRPCServer(tourService)
		pb.RegisterTourServiceServer(grpcServer, tourGRPCServer)

		lis, err := net.Listen("tcp", ":50053")
		if err != nil {
			log.Fatalf("Failed to listen on port 50053: %v", err)
		}

		fmt.Println("Tour gRPC server running on port 50053")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	fmt.Println("Tour service running on internal port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
