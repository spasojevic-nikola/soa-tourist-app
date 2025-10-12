package grpc

import (
	"context"
	"log"
	"time"

	pb "tour-service/gen/pb-go/tour"
	"tour-service/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TourGRPCServer implementira gRPC servis definisan u tour.proto
type TourGRPCServer struct {
	pb.UnimplementedTourServiceServer
	TourService   *service.TourService
	ReviewService *service.ReviewService
}

// NewTourGRPCServer kreira novu instancu gRPC servera
func NewTourGRPCServer(tourService *service.TourService, reviewService *service.ReviewService) *TourGRPCServer {
	return &TourGRPCServer{
		TourService:   tourService,
		ReviewService: reviewService,
	}
}

// GetPublishedTours implementira gRPC metodu za preuzimanje objavljenih tura
func (s *TourGRPCServer) GetPublishedTours(ctx context.Context, req *pb.GetPublishedToursRequest) (*pb.GetPublishedToursResponse, error) {
	log.Println("gRPC: GetPublishedTours pozvan")

	// Pozivamo servis za dobijanje objavljenih tura
	tours, err := s.TourService.GetAllPublishedTours()
	if err != nil {
		log.Printf("Greška pri preuzimanju objavljenih tura: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to retrieve published tours: %v", err)
	}

	// Konvertujemo Go modele u Protobuf poruke
	var pbTours []*pb.Tour
	for _, tour := range tours {
		pbTour := &pb.Tour{
			Id:          uint32(tour.ID),
			AuthorId:    uint32(tour.AuthorID),
			Name:        tour.Name,
			Description: tour.Description,
			Difficulty:  string(tour.Difficulty),
			Tags:        tour.Tags,
			Status:      string(tour.Status),
			Price:       tour.Price,
			IsDeleted:   tour.IsDeleted,
			CreatedAt:   tour.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   tour.UpdatedAt.Format(time.RFC3339),
		}

		// Konvertujemo KeyPoints
		var pbKeyPoints []*pb.KeyPoint
		for _, kp := range tour.KeyPoints {
			pbKeyPoint := &pb.KeyPoint{
				Id:          uint32(kp.ID),
				TourId:      uint32(kp.TourID),
				Name:        kp.Name,
				Description: kp.Description,
				Latitude:    kp.Latitude,
				Longitude:   kp.Longitude,
				Image:       kp.Image,
				Order:       uint32(kp.Order),
				CreatedAt:   kp.CreatedAt.Format(time.RFC3339),
				UpdatedAt:   kp.UpdatedAt.Format(time.RFC3339),
			}
			pbKeyPoints = append(pbKeyPoints, pbKeyPoint)
		}
		pbTour.KeyPoints = pbKeyPoints

		pbTours = append(pbTours, pbTour)
	}

	response := &pb.GetPublishedToursResponse{
		Tours:      pbTours,
		TotalCount: int32(len(pbTours)),
	}

	log.Printf("gRPC: Vraćam %d objavljenih tura", len(pbTours))
	return response, nil
}

// GetMyTours implementira gRPC metodu za preuzimanje tura autora
func (s *TourGRPCServer) GetMyTours(ctx context.Context, req *pb.GetMyToursRequest) (*pb.GetMyToursResponse, error) {
	log.Printf("gRPC: GetMyTours pozvan za author_id: %d", req.AuthorId)

	// Pozivamo servis za dobijanje tura autora
	tours, err := s.TourService.GetToursByAuthor(uint(req.AuthorId))
	if err != nil {
		log.Printf("Greška pri preuzimanju tura autora: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to retrieve author tours: %v", err)
	}

	// Konvertujemo Go modele u Protobuf poruke
	var pbTours []*pb.Tour
	for _, tour := range tours {
		pbTour := &pb.Tour{
			Id:          uint32(tour.ID),
			AuthorId:    uint32(tour.AuthorID),
			Name:        tour.Name,
			Description: tour.Description,
			Difficulty:  string(tour.Difficulty),
			Tags:        tour.Tags,
			Status:      string(tour.Status),
			Price:       tour.Price,
			IsDeleted:   tour.IsDeleted,
			CreatedAt:   tour.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   tour.UpdatedAt.Format(time.RFC3339),
		}

		// Konvertujemo KeyPoints
		var pbKeyPoints []*pb.KeyPoint
		for _, kp := range tour.KeyPoints {
			pbKeyPoint := &pb.KeyPoint{
				Id:          uint32(kp.ID),
				TourId:      uint32(kp.TourID),
				Name:        kp.Name,
				Description: kp.Description,
				Latitude:    kp.Latitude,
				Longitude:   kp.Longitude,
				Image:       kp.Image,
				Order:       uint32(kp.Order),
				CreatedAt:   kp.CreatedAt.Format(time.RFC3339),
				UpdatedAt:   kp.UpdatedAt.Format(time.RFC3339),
			}
			pbKeyPoints = append(pbKeyPoints, pbKeyPoint)
		}
		pbTour.KeyPoints = pbKeyPoints

		pbTours = append(pbTours, pbTour)
	}

	response := &pb.GetMyToursResponse{
		Tours:      pbTours,
		TotalCount: int32(len(pbTours)),
	}

	log.Printf("gRPC: Vraćam %d tura za autora %d", len(pbTours), req.AuthorId)
	return response, nil
}

// GetReviewsByTour implementira gRPC metodu za preuzimanje recenzija ture
func (s *TourGRPCServer) GetReviewsByTour(ctx context.Context, req *pb.GetReviewsByTourRequest) (*pb.GetReviewsByTourResponse, error) {
	log.Printf("gRPC: GetReviewsByTour pozvan za tour_id: %d", req.TourId)

	// Pozivamo servis za dobijanje recenzija
	reviews, err := s.ReviewService.GetReviewsByTourID(uint(req.TourId))
	if err != nil {
		log.Printf("Greška pri preuzimanju recenzija: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to retrieve reviews: %v", err)
	}

	// Konvertujemo Go modele u Protobuf poruke
	var pbReviews []*pb.Review
	for _, review := range reviews {
		pbReview := &pb.Review{
			Id:              uint32(review.ID),
			TourId:          uint32(review.TourID),
			TouristId:       uint32(review.TouristID),
			TouristUsername: review.TouristUsername,
			Rating:          int32(review.Rating),
			Comment:         review.Comment,
			VisitDate:       review.VisitDate.Format(time.RFC3339),
			Images:          review.Images,
			CreatedAt:       review.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       review.UpdatedAt.Format(time.RFC3339),
		}
		pbReviews = append(pbReviews, pbReview)
	}

	response := &pb.GetReviewsByTourResponse{
		Reviews:    pbReviews,
		TotalCount: int32(len(pbReviews)),
	}

	log.Printf("gRPC: Vraćam %d recenzija za turu %d", len(pbReviews), req.TourId)
	return response, nil
}
