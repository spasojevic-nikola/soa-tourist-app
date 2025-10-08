package grpc

import (
	"context"
	"log"
	"net"

	"blog-service/gen/pb-go"
	"blog-service/internal/service"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	blog.UnimplementedBlogServiceServer
	blogService *service.BlogService
}

func NewgRPCServer(blogService *service.BlogService) *gRPCServer {
	return &gRPCServer{
		blogService: blogService,
	}
}

func (s *gRPCServer) GetAllBlogs(ctx context.Context, req *blog.GetBlogsRequest) (*blog.GetBlogsResponse, error) {
	// Pozovi postojeću metodu sa context-om
	blogs, err := s.blogService.GetAllBlogs(ctx)
	if err != nil {
		return nil, err
	}

	// Konvertuj domain blogove u gRPC blogove
	var grpcBlogs []*blog.Blog
	for _, b := range blogs {
		grpcBlogs = append(grpcBlogs, &blog.Blog{
			Id:          b.ID.Hex(), // Konvertuj ObjectID u string
			Title:       b.Title,
			Description: b.Content, 
			Author:      string(b.AuthorID), // Konvertuj uint u string
			CreatedAt:   b.CreatedAt.Format("2006-01-02 15:04:05"),
			LikesCount:    int32(len(b.Likes)),
			CommentsCount: int32(len(b.Comments)),
		})
	}

	return &blog.GetBlogsResponse{
		Blogs:      grpcBlogs,
		TotalCount: int32(len(grpcBlogs)),
	}, nil
}

func StartGRPCServer(blogService *service.BlogService, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := NewgRPCServer(blogService)
	blog.RegisterBlogServiceServer(grpcServer, server)

	log.Printf("Blog gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}