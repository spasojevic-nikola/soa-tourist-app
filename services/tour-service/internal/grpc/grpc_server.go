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
	TourService *service.TourService
}

// NewTourGRPCServer kreira novu instancu gRPC servera
func NewTourGRPCServer(tourService *service.TourService) *TourGRPCServer {
	return &TourGRPCServer{
		TourService: tourService,
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
