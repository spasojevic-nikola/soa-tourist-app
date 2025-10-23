package main

import (
	"fmt"
	"net/http"
	"time"
	"soa-tourist-app/follower-service/internal/api"
	"soa-tourist-app/follower-service/internal/database"
	"soa-tourist-app/follower-service/internal/repository"
	"soa-tourist-app/follower-service/internal/service"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Konfigurisanje logrus-a za JSON format (isti kao blog service)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(false)
}

func main() {

	log.WithFields(log.Fields{
		"service": "follower-service",
		"port":    "8080",
	}).Info("Starting follower service")

	// Inicijalizacija
	driver := database.InitDB()
	// defer driver.Close(context.Background()) // defer se sada ne koristi jer server radi non-stop

	repo := repository.NewFollowerRepository(driver)
	followerService := service.NewFollowerService(repo)
	handler := api.NewHandler(followerService)

	// Ruter
	r := mux.NewRouter()
	
	// Kreiramo pod-ruter za rute koje zahtevaju autorizaciju
	protectedRoutes := r.PathPrefix("/api/followers").Subrouter()
	protectedRoutes.Use(api.AuthMiddleware) // Primenjujemo middleware na sve rute u ovom pod-ruteru

	// Definišemo rute na pod-ruteru
	protectedRoutes.HandleFunc("/follow/{id:[0-9]+}", handler.Follow).Methods("POST")
	protectedRoutes.HandleFunc("/unfollow/{id:[0-9]+}", handler.Unfollow).Methods("DELETE")
	protectedRoutes.HandleFunc("/check-follow/{id:[0-9]+}", handler.CheckFollow).Methods("GET")
	protectedRoutes.HandleFunc("/recommendations", handler.GetRecommendations).Methods("GET")
	protectedRoutes.HandleFunc("/following", handler.GetFollowingIDs).Methods("GET")
	
	// Podesavanje CORS-a
	//corsHandler := handlers.CORS(
	//	handlers.AllowedOrigins([]string{"http://localhost:4200"}), // Promenite ako je potrebno
	//	handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
	//	handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-User-ID"}),
	//	)(r)

		log.WithFields(log.Fields{
			"port": "8080",
		}).Info("Follower service is ready to accept connections")

	// Pokretanje servera
	fmt.Println("Follower service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrapper za response writer da bismo uhvatili status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapped, r)

		log.WithFields(log.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      wrapped.statusCode,
			"duration_ms": time.Since(start).Milliseconds(),
			"remote_addr": r.RemoteAddr,
		}).Info("HTTP request")
	})
}

// Custom response writer za uhvatanje status koda
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}