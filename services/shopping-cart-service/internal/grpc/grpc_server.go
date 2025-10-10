package grpc

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"shopping-cart-service/gen/pb-go"   
	"shopping-cart-service/internal/service"
)

// gRPCServer implementira ShoppingCartServiceServer
type gRPCServer struct {
	shopping_cart.UnimplementedShoppingCartServiceServer
	cartService *service.CartService
}

// NewgRPCServer kreira novi gRPC server
func NewgRPCServer(cartService *service.CartService) *gRPCServer {
	return &gRPCServer{
		cartService: cartService,
	}
}

// VerifyPurchase implementira RPC metodu
func (s *gRPCServer) VerifyPurchase(ctx context.Context, req *shopping_cart.VerifyPurchaseRequest) (*shopping_cart.VerifyPurchaseResponse, error) {
	// Konvertuj tour_id iz string u uint
	// Ovdje koristimo postojeći CartService
	hasPurchased, err := s.cartService.HasPurchaseToken(ctx, uint(req.GetTouristId()), req.GetTourId())
	if err != nil {
		return &shopping_cart.VerifyPurchaseResponse{
			HasPurchased: false,
			Message:      "Error checking purchase: " + err.Error(),
		}, nil
	}

	return &shopping_cart.VerifyPurchaseResponse{
		HasPurchased: hasPurchased,
		Message:      "Purchase verification completed",
	}, nil
}

// Start pokreće gRPC server
func StartGRPCServer(cartService *service.CartService, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := NewgRPCServer(cartService)
	shopping_cart.RegisterShoppingCartServiceServer(grpcServer, server)

	log.Printf("gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}