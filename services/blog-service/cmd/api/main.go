package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"blog-service/internal/api"
	"blog-service/internal/database"
	"blog-service/internal/repository" 
	"blog-service/internal/service" 	

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	mongoDB := database.InitDB() 

	blogRepo := repository.NewBlogRepository(mongoDB)

	blogService := service.NewBlogService(blogRepo)

	blogHandler := api.NewHandler(blogService) 

	r := mux.NewRouter()

	// Definisemo CORS opcije
	corsOpts := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	apiV1 := r.PathPrefix("/api/v1/blogs").Subrouter()

	apiV1.HandleFunc("", api.AuthMiddleware(blogHandler.CreateBlog)).Methods("POST")
	apiV1.HandleFunc("/{id}/comments", api.AuthMiddleware(blogHandler.AddComment)).Methods("POST")
	apiV1.HandleFunc("/{id}/like", api.AuthMiddleware(blogHandler.ToggleLike)).Methods("POST")
	apiV1.HandleFunc("", blogHandler.GetAllBlogs).Methods("GET")
	apiV1.HandleFunc("/{id}", blogHandler.GetBlogByID).Methods("GET")


	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	fmt.Println("Blog service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", corsOpts(r)))
}