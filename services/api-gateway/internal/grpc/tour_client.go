package grpc

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"api-gateway/gen/pb-go/tour"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TourClient struct {
	client tour.TourServiceClient
	conn   *grpc.ClientConn
}

func NewTourClient(serverAddr string) (*TourClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := tour.NewTourServiceClient(conn)

	return &TourClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *TourClient) GetPublishedToursHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("API Gateway: Primljen REST zahtev za objavljene ture")

	resp, err := c.client.GetPublishedTours(context.Background(), &tour.GetPublishedToursRequest{})
	if err != nil {
		log.Printf("gRPC greska: %v", err)
		http.Error(w, "Tour service unavailable", http.StatusServiceUnavailable)
		return
	}

	log.Printf("API Gateway: Primljeno %d tura od Tour servisa preko gRPC", resp.TotalCount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetMyToursHandler handler za REST endpoint koji poziva gRPC servis za ture autora
func (c *TourClient) GetMyToursHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("API Gateway: Primljen REST zahtev za moje ture")

	// Izvlačimo userID iz headera (postavljeno od AuthMiddleware)
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		log.Println("User ID nije pronađen u zahtevu")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Printf("Nevažeći User ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("API Gateway: Pozivam Tour Service za author_id: %d", userID)

	// Pozovi gRPC Tour servis
	resp, err := c.client.GetMyTours(context.Background(), &tour.GetMyToursRequest{
		AuthorId: uint32(userID),
	})
	if err != nil {
		log.Printf("gRPC greska: %v", err)
		http.Error(w, "Tour service unavailable", http.StatusServiceUnavailable)
		return
	}

	log.Printf("API Gateway: Primljeno %d tura od Tour servisa preko gRPC", resp.TotalCount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetReviewsByTourHandler handler za REST endpoint koji poziva gRPC servis za recenzije ture
func (c *TourClient) GetReviewsByTourHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("API Gateway: Primljen REST zahtev za recenzije ture")

	// Izvlačimo tourID iz URL-a (očekujemo format /tour/{tourId}/reviews)
	pathParts := strings.Split(r.URL.Path, "/")
	var tourIDStr string
	for i, part := range pathParts {
		if part == "tour" && i+1 < len(pathParts) {
			tourIDStr = pathParts[i+1]
			break
		}
	}

	if tourIDStr == "" {
		log.Println("Tour ID nije pronađen u URL-u")
		http.Error(w, "Tour ID is required", http.StatusBadRequest)
		return
	}

	tourID, err := strconv.ParseUint(tourIDStr, 10, 32)
	if err != nil {
		log.Printf("Nevažeći Tour ID: %v", err)
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	log.Printf("API Gateway: Pozivam Tour Service za tour_id: %d", tourID)

	// Pozovi gRPC Tour servis
	resp, err := c.client.GetReviewsByTour(context.Background(), &tour.GetReviewsByTourRequest{
		TourId: uint32(tourID),
	})
	if err != nil {
		log.Printf("gRPC greska: %v", err)
		http.Error(w, "Tour service unavailable", http.StatusServiceUnavailable)
		return
	}

	log.Printf("API Gateway: Primljeno %d recenzija od Tour servisa preko gRPC", resp.TotalCount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *TourClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
