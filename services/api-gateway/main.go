package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"api-gateway/internal/middleware"
)

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

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-User-Username, X-User-Role")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Handler koji rutira zahteve
func router(w http.ResponseWriter, r *http.Request) {
	// Normalize trailing slashes (remove trailing / except for root)
	if len(r.URL.Path) > 1 && strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	}
	
	path := r.URL.Path

	// =========================== BLOG SERVICE ===========================
	blogProxy := newReverseProxy("http://blog-service:8081", "")

	// =========================== AUTH SERVICE ===========================
	if strings.HasPrefix(path, "/api/v1/auth") {
		log.Printf("Routing Auth: %s %s", r.Method, path)
		proxy := newReverseProxy("http://auth-service:8084", "")
		proxy.ServeHTTP(w, r)
		return
	}

	// =========================== BLOG SERVICE ===========================
	if strings.HasPrefix(path, "/api/v1/blogs") {
		// Protected blog routes
		if (r.Method == "POST" && path == "/api/v1/blogs") ||
			(r.Method == "POST" && strings.HasSuffix(path, "/comments")) ||
			(r.Method == "POST" && strings.HasSuffix(path, "/like")) ||
			(r.Method == "PUT" && strings.HasPrefix(path, "/api/v1/blogs/") && !strings.Contains(path, "/comments")) ||
			(r.Method == "PUT" && strings.Contains(path, "/comments/")) {
			log.Printf("Routing %s to Blog Service (PROTECTED): %s", r.Method, path)
			middleware.JWTAuthMiddleware(blogProxy).ServeHTTP(w, r)
			return
		}
		// Public blog routes - use REST instead of gRPC to support feed logic
		log.Printf("Routing %s to Blog Service [PUBLIC]: %s", r.Method, path)
		middleware.OptionalJWTMiddleware(blogProxy).ServeHTTP(w, r)
		return
	}

	// =========================== STAKEHOLDERS SERVICE ===========================
	if strings.HasPrefix(path, "/api/v1/user") || 
	   strings.HasPrefix(path, "/api/v1/profile") ||
	   strings.HasPrefix(path, "/api/v1/admin/users") ||
	   strings.HasPrefix(path, "/api/v1/users") {
		log.Printf("Routing Stakeholders: %s %s", r.Method, path)
		proxy := newReverseProxy("http://stakeholders-service:8080", "")
		proxy.ServeHTTP(w, r)
		return
	}

	// =========================== FOLLOWER SERVICE ===========================
	if strings.HasPrefix(path, "/api/followers") {
		log.Printf("Routing Follower: %s %s", r.Method, path)
		proxy := newReverseProxy("http://follower-service:8080", "")
		middleware.JWTAuthMiddleware(proxy).ServeHTTP(w, r)
		return
	}

	// =========================== TOUR SERVICE ===========================
	if strings.HasPrefix(path, "/api/v1/tours") {
		log.Printf("Routing Tour: %s %s", r.Method, path)
		proxy := newReverseProxy("http://tour-service:8080", "")
		proxy.ServeHTTP(w, r)
		return
	}

	// =========================== SHOPPING CART SERVICE ===========================
	if strings.HasPrefix(path, "/api/v1/cart") {
		log.Printf("Routing Shopping Cart: %s %s", r.Method, path)
		proxy := newReverseProxy("http://shopping-cart-service:8081", "")
		proxy.ServeHTTP(w, r)
		return
	}

	// =========================== HEALTH CHECKS ===========================
	if path == "/health" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "api-gateway"}`))
		return
	}

	// Default: Not Found
	log.Printf("Route not found: %s %s", r.Method, path)
	http.Error(w, "Not Found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/", corsMiddleware(router))
	log.Println("API Gateway running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
