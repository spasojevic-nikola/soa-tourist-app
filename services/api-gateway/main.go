package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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

	case r.Method == "GET" && path == "/api/v1/blogs":
		log.Printf("Routing GET %s to Blog Service (Get All Blogs)", path)
		blogProxy.ServeHTTP(w, r)

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
	case strings.HasPrefix(path, "/tour"):
		log.Printf("Routing Tour: %s", path)
		proxy := newReverseProxy("http://tour-service:8080", "/tour")
		proxy.ServeHTTP(w, r)

	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/", router)
	log.Println("API Gateway running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
