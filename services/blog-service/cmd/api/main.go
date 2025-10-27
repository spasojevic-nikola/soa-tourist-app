package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"blog-service/internal/api"
	"blog-service/internal/database"
	"blog-service/internal/grpc"
	"blog-service/internal/repository"
	"blog-service/internal/service"
	"blog-service/internal/tracing"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func init() {
	// Konfigurisanje logrus-a za JSON format
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(false)
}

func main() {
	// Initialize tracing
	cleanup, err := tracing.InitTracing("blog-service", "1.0.0")
	if err != nil {
		log.WithError(err).Warn("Failed to initialize tracing, continuing without tracing")
	} else {
		// Setup graceful shutdown for tracing
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			if err := cleanup(); err != nil {
				log.WithError(err).Error("Failed to shutdown tracer")
			}
			os.Exit(0)
		}()
	}

	log.WithFields(log.Fields{
		"service": "blog-service",
		"port":    "8081",
	}).Info("Starting blog service with tracing enabled")

	mongoDB := database.InitDB()

	blogRepo := repository.NewBlogRepository(mongoDB)

	blogService := service.NewBlogService(blogRepo)

	blogHandler := api.NewHandler(blogService)

	// pokreni gRPC server u pozadini
	go func() {
		log.WithFields(log.Fields{
			"port": "50052",
		}).Info("Starting Blog gRPC server")
		grpc.StartGRPCServer(blogService, "50052")
	}()

	r := mux.NewRouter()

	// Add OpenTelemetry middleware for HTTP tracing
	r.Use(otelmux.Middleware("blog-service"))

	// CORS removed - requests now go through API gateway which handles CORS

	apiV1 := r.PathPrefix("/api/v1/blogs").Subrouter()

	// API Gateway sada radi JWT validaciju, mi samo čitamo X-User-* headere
	apiV1.HandleFunc("", blogHandler.CreateBlog).Methods("POST")
	apiV1.HandleFunc("/{id}/comments", blogHandler.AddComment).Methods("POST")
	apiV1.HandleFunc("/{id}/like", blogHandler.ToggleLike).Methods("POST")
	apiV1.HandleFunc("", blogHandler.GetAllBlogs).Methods("GET")
	apiV1.HandleFunc("/{id}", blogHandler.GetBlogByID).Methods("GET")
	apiV1.HandleFunc("/{id}", blogHandler.UpdateBlog).Methods("PUT")
	apiV1.HandleFunc("/{id}/comments/{commentId}", blogHandler.UpdateComment).Methods("PUT")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	log.WithFields(log.Fields{
		"port": "8081",
	}).Info("Blog service is ready to accept connections")

	fmt.Println("Blog service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
