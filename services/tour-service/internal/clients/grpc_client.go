package clients

import (
	"context"
	"log"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"tour-service/gen/pb-go" // Generisani kod
)

type GRPCPurchaseChecker struct {
	client shopping_cart.ShoppingCartServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCPurchaseChecker(serverAddr string) (*GRPCPurchaseChecker, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := shopping_cart.NewShoppingCartServiceClient(conn)
	
	return &GRPCPurchaseChecker{
		client: client,
		conn:   conn,
	}, nil
}

func (g *GRPCPurchaseChecker) HasPurchasedTour(touristID uint, tourID uint) (bool, error) {
	// Konvertuj uint tourID u string za gRPC (MongoDB koristi string ID-eve)
	tourIDStr := strconv.FormatUint(uint64(tourID), 10)
	
	resp, err := g.client.VerifyPurchase(context.Background(), &shopping_cart.VerifyPurchaseRequest{
		TouristId: uint32(touristID),
		TourId:    tourIDStr,
	})
	
	if err != nil {
		log.Printf("gRPC error: %v", err)
		return false, err
	}

	return resp.HasPurchased, nil
}

func (g *GRPCPurchaseChecker) Close() {
	if g.conn != nil {
		g.conn.Close()
	}
}