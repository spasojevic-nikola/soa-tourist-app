package grpc

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"api-gateway/gen/pb-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BlogClient struct {
	client blog.BlogServiceClient
	conn   *grpc.ClientConn
}

func NewBlogClient(serverAddr string) (*BlogClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := blog.NewBlogServiceClient(conn)
	
	return &BlogClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *BlogClient) GetAllBlogsHandler(w http.ResponseWriter, r *http.Request) {
	// Pozovi gRPC servis
	resp, err := c.client.GetAllBlogs(context.Background(), &blog.GetBlogsRequest{})
	if err != nil {
		log.Printf("gRPC error: %v", err)
		http.Error(w, "Blog service unavailable", http.StatusServiceUnavailable)
		return
	}

	// Konvertuj gRPC odgovor u JSON i pošalji klijentu
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *BlogClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}