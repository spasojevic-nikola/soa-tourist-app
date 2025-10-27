package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	blogpb "blog-service/proto"
	"blog-service/internal/models"
	"blog-service/internal/service"
	"google.golang.org/grpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

type gRPCServer struct {
	blogpb.UnimplementedBlogServiceServer
	blogService *service.BlogService
}

func NewgRPCServer(blogService *service.BlogService) *gRPCServer {
	return &gRPCServer{
		blogService: blogService,
	}
}

func (s *gRPCServer) GetAllBlogs(ctx context.Context, req *blogpb.GetBlogsRequest) (*blogpb.GetBlogsResponse, error) {
	// Konvertuj uint32 user_id u uint i pozovi GetFeedForUser za feed logiku
	userID := uint(req.UserId)
	
	var blogs []*models.Blog
	
	if userID == 0 {
		// Anonymous user - get all blogs
		allBlogs, err := s.blogService.GetAllBlogs(ctx)
		if err != nil {
			return nil, err
		}
		blogs = make([]*models.Blog, len(allBlogs))
		for i := range allBlogs {
			blogs[i] = &allBlogs[i]
		}
	} else {
		// Authenticated user - get feed for user
		feedBlogs, err := s.blogService.GetFeedForUser(ctx, userID)
		if err != nil {
			return nil, err
		}
		blogs = make([]*models.Blog, len(feedBlogs))
		for i := range feedBlogs {
			blogs[i] = &feedBlogs[i]
		}
	}

	// Konvertuj domain blogove u gRPC blogove
	var grpcBlogs []*blogpb.Blog
	for _, b := range blogs {
		grpcBlogs = append(grpcBlogs, &blogpb.Blog{
			Id:          b.ID.Hex(), // Konvertuj ObjectID u string
			Title:       b.Title,
			Description: b.Content, 
			Author:      fmt.Sprintf("%d", b.AuthorID), // Konvertuj uint u string
			CreatedAt:   b.CreatedAt.Format("2006-01-02 15:04:05"),
			LikesCount:    int32(len(b.Likes)),
			CommentsCount: int32(len(b.Comments)),
		})
	}

	return &blogpb.GetBlogsResponse{
		Blogs:      grpcBlogs,
		TotalCount: int32(len(grpcBlogs)),
	}, nil
}

func StartGRPCServer(blogService *service.BlogService, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server with OpenTelemetry interceptors for tracing
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	server := NewgRPCServer(blogService)
	blogpb.RegisterBlogServiceServer(grpcServer, server)

	log.Printf("Blog gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}