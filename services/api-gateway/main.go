package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"api-gateway/internal/grpc"
)

var blogClient *grpc.BlogClient
var tourClient *grpc.TourClient

func init() {
	// DODAJ: Inicijalizacija gRPC klijenta pri pokretanju
	var err error
	blogClient, err = grpc.NewBlogClient("blog-service:50052")
	if err != nil {
		log.Fatalf("Failed to create gRPC Blog client: %v", err)
	}
	log.Println("gRPC Blog client initialized")

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

// Handler koji rutira zahteve
func router(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// =========================== BLOG SERVICE ===========================
	blogProxy := newReverseProxy("http://blog-service:8081", "")

	switch {
	case r.Method == "POST" && path == "/api/v1/blogs":
		log.Printf("Routing POST %s to Blog Service (Create Blog)", path)
		blogProxy.ServeHTTP(w, r)

	case r.Method == "POST" && strings.HasSuffix(path, "/comments"):
		log.Printf("Routing POST %s to Blog Service (Add Comment)", path)
		blogProxy.ServeHTTP(w, r)

	case r.Method == "POST" && strings.HasSuffix(path, "/like"):
		log.Printf("Routing POST %s to Blog Service (Toggle Like)", path)
		blogProxy.ServeHTTP(w, r)

	// promijenjeno iz rest zahtjeva u rpc poziv
	case r.Method == "GET" && path == "/api/v1/blogs":
		log.Printf("Routing GET %s to Blog Service via gRPC", path)
		blogClient.GetAllBlogsHandler(w, r)

	case r.Method == "GET" && strings.HasPrefix(path, "/api/v1/blogs/"):
		log.Printf("Routing GET %s to Blog Service (Get Blog By ID)", path)
		blogProxy.ServeHTTP(w, r)

	case r.Method == "PUT" && strings.HasPrefix(path, "/api/v1/blogs/") && !strings.Contains(path, "/comments"):
		log.Printf("Routing PUT %s to Blog Service (Update Blog)", path)
		blogProxy.ServeHTTP(w, r)

	case r.Method == "PUT" && strings.Contains(path, "/comments/"):
		log.Printf("Routing PUT %s to Blog Service (Update Comment)", path)
		blogProxy.ServeHTTP(w, r)

	// ====================== AUTH SERVICE ======================
	case strings.HasPrefix(path, "/api/v1/auth"):
		log.Printf("Routing Auth: %s", path)
		proxy := newReverseProxy("http://auth-service:8084", "/")
		proxy.ServeHTTP(w, r)

	// ====================== STAKEHOLDERS SERVICE ======================
	case strings.HasPrefix(path, "/stakeholders"):
		log.Printf("Routing Stakeholders: %s", path)
		proxy := newReverseProxy("http://stakeholders-service:8080", "/stakeholders")
		proxy.ServeHTTP(w, r)

	// ====================== FOLLOWER SERVICE ======================
	case strings.HasPrefix(path, "/follower"):
		log.Printf("Routing Follower: %s", path)
		proxy := newReverseProxy("http://follower-service:8080", "/follower")
		proxy.ServeHTTP(w, r)

	// ====================== TOUR SERVICE ======================
	case r.Method == "GET" && path == "/api/v1/tours/published":
		log.Printf("Routing GET %s to Tour Service via gRPC", path)
		tourClient.GetPublishedToursHandler(w, r)

	case r.Method == "GET" && path == "/api/v1/tours":
		log.Printf("Routing GET %s to Tour Service via gRPC (My Tours)", path)
		tourClient.GetMyToursHandler(w, r)

	case r.Method == "GET" && strings.Contains(path, "/reviews") && !strings.Contains(path, "/reviews/"):
		log.Printf("Routing GET %s to Tour Service via gRPC (Reviews)", path)
		tourClient.GetReviewsByTourHandler(w, r)

	case strings.HasPrefix(path, "/tour"):
		log.Printf("Routing Tour: %s", path)
		proxy := newReverseProxy("http://tour-service:8080", "/tour")
		proxy.ServeHTTP(w, r)

	default:
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

	http.HandleFunc("/", router)
	log.Println("API Gateway running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
