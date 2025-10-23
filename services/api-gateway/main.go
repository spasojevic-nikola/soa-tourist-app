package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"api-gateway/internal/grpc"
	"api-gateway/internal/middleware"
    // Novi uvoz za CORS
    "github.com/gorilla/handlers" 
)

var blogClient *grpc.BlogClient
var tourClient *grpc.TourClient

func init() {
	var err error
	
	// Inicijalizacija gRPC klijenta za Blog Service
	blogClient, err = grpc.NewBlogClient("blog-service:50052")
	if err != nil {
		log.Fatalf("Failed to create gRPC Blog client: %v", err)
	}
	log.Println("gRPC Blog client initialized")

	// Inicijalizacija gRPC klijenta za Tour Service
	tourClient, err = grpc.NewTourClient("tour-service:50053")
	if err != nil {
		log.Fatalf("Failed to create gRPC Tour client: %v", err)
	}
	log.Println("gRPC Tour client initialized")
}

func newReverseProxy(targetURL string, pathPrefix string) *httputil.ReverseProxy {
	url, _ := url.Parse(targetURL)
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host

		if pathPrefix != "" && pathPrefix != "/" {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, pathPrefix)
		}

		if req.URL.Path == "" {
			req.URL.Path = "/"
		}

		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = url.Host
	}
	return proxy
}

// Handler koji rutira zahteve (isto kao pre)
func router(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Printf("Received request: %s %s", r.Method, path)

	blogProxy := newReverseProxy("http://blog-service:8081", "/")


	switch {
	// ====================== AUTH SERVICE ======================
	case strings.HasPrefix(path, "/api/v1/auth"):
		log.Printf("Routing Auth: %s", path)
		proxy := newReverseProxy("http://auth-service:8084", "/") 
		proxy.ServeHTTP(w, r)


	// ====================== STAKEHOLDERS SERVICE ======================
	case strings.HasPrefix(path, "/api/v1/user"), 
		strings.HasPrefix(path, "/api/v1/users"), 
		strings.HasPrefix(path, "/api/v1/profile"), 
		strings.HasPrefix(path, "/api/v1/admin/users"):
		
		log.Printf("Routing Stakeholders: %s", path)
		proxy := newReverseProxy("http://stakeholders-service:8080", "/api/v1") 
		middleware.JWTAuthMiddleware(proxy).ServeHTTP(w, r)


	// =========================== BLOG SERVICE ===========================
	case r.Method == "POST" && path == "/api/v1/blogs":
		log.Printf("Routing POST %s to Blog Service (Create Blog) [AUTH REQUIRED]", path)
		middleware.JWTAuthMiddleware(blogProxy).ServeHTTP(w, r)
    // ... (ostale blog rute ostaju nepromenjene)

	case r.Method == "POST" && strings.HasSuffix(path, "/comments"):
		log.Printf("Routing POST %s to Blog Service (Add Comment) [AUTH REQUIRED]", path)
		middleware.JWTAuthMiddleware(blogProxy).ServeHTTP(w, r)

	case r.Method == "POST" && strings.HasSuffix(path, "/like"):
		log.Printf("Routing POST %s to Blog Service (Toggle Like) [AUTH REQUIRED]", path)
		middleware.JWTAuthMiddleware(blogProxy).ServeHTTP(w, r)

	case r.Method == "PUT" && strings.HasPrefix(path, "/api/v1/blogs/") && !strings.Contains(path, "/comments"):
		log.Printf("Routing PUT %s to Blog Service (Update Blog) [AUTH REQUIRED]", path)
		middleware.JWTAuthMiddleware(blogProxy).ServeHTTP(w, r)

	case r.Method == "PUT" && strings.Contains(path, "/comments/"):
		log.Printf("Routing PUT %s to Blog Service (Update Comment) [AUTH REQUIRED]", path)
		middleware.JWTAuthMiddleware(blogProxy).ServeHTTP(w, r)

	case r.Method == "GET" && path == "/api/v1/blogs":
		log.Printf("Routing GET %s to Blog Service via gRPC [PUBLIC]", path)
		middleware.OptionalJWTMiddleware(http.HandlerFunc(blogClient.GetAllBlogsHandler)).ServeHTTP(w, r)

	case r.Method == "GET" && strings.HasPrefix(path, "/api/v1/blogs/"):
		log.Printf("Routing GET %s to Blog Service (Get Blog By ID) [PUBLIC]", path)
		middleware.OptionalJWTMiddleware(blogProxy).ServeHTTP(w, r)


	// ====================== FOLLOWER SERVICE ======================
	case strings.HasPrefix(path, "/api/followers"):
		log.Printf("Routing Follower: %s", path)
		proxy := newReverseProxy("http://follower-service:8080", "/")
		middleware.JWTAuthMiddleware(proxy).ServeHTTP(w, r)
	
	
	// ====================== SHOPPING CART SERVICE ======================
	case strings.HasPrefix(path, "/api/v1/cart"):
		log.Printf("Routing Cart: %s", path)
		proxy := newReverseProxy("http://shopping-cart-service:8081", "/")
		middleware.JWTAuthMiddleware(proxy).ServeHTTP(w, r)


	// ====================== TOUR SERVICE ======================
	case r.Method == "GET" && path == "/api/v1/tours/published":
		log.Printf("Routing GET %s to Tour Service via gRPC", path)
		tourClient.GetPublishedToursHandler(w, r)

	case r.Method == "GET" && path == "/api/v1/tours":
		log.Printf("Routing GET %s to Tour Service via gRPC (My Tours)", path)
		tourClient.GetMyToursHandler(w, r)
		
	case strings.HasPrefix(path, "/api/v1/tours"):
		log.Printf("Routing Tour (REST): %s", path)
		proxy := newReverseProxy("http://tour-service:8080", "/") 
		middleware.JWTAuthMiddleware(proxy).ServeHTTP(w, r)


	default:
		// Dodaj Health Check rutu ovde, ako je nisi stavio negde drugde
		if path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func main() {
	defer func() {
		if blogClient != nil {
			blogClient.Close()
		}
		if tourClient != nil {
			tourClient.Close()
		}
	}()

    // 1. Definicija CORS opcija
    corsOpts := handlers.CORS(
        // Ovo mora da odgovara Angular adresi
        handlers.AllowedOrigins([]string{"http://localhost:4200"}), 
        // Dozvoli sve metode koje koristiš
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
        // Ključno: Dozvoli headere koje šalješ, posebno Authorization
        handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-User-ID", "X-User-Username", "X-User-Role"}),
    )

    // 2. Umotaj router u CORS middleware i pokreni server
	log.Println("API Gateway running on port 8080...")
    // Koristimo corsOpts(http.HandlerFunc(router)) umesto nil
	if err := http.ListenAndServe(":8080", corsOpts(http.HandlerFunc(router))); err != nil {
		log.Fatal(err)
	}
}