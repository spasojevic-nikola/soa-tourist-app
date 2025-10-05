package api

import (
	"encoding/json"
	"net/http"

	"blog-service/internal/dto"
	"blog-service/internal/service"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	log "github.com/sirupsen/logrus"
)


// Handler sadrži referencu na BlogService
type Handler struct {
	Service *service.BlogService
}

// NewHandler kreira novu instancu Handler-a
func NewHandler(service *service.BlogService) *Handler {
	return &Handler{Service: service}
}

// CreateBlog endpoint za kreiranje bloga
func (h *Handler) CreateBlog(w http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{
		"endpoint": "/api/v1/blogs",
		"method":   "POST",
		"ip":       r.RemoteAddr,
	}).Info("Create blog request received")

	var req dto.CreateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	authorID := r.Context().Value("userID").(uint)

	blog, err := h.Service.CreateBlog(r.Context(), req, authorID)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": "/api/v1/blogs",
			"authorID": authorID,
			"error":    err.Error(),
		}).Error("Failed to create blog")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.WithFields(log.Fields{
		"endpoint": "/api/v1/blogs",
		"authorID": authorID,
		"blogID":   blog.ID.Hex(),
		"title":    req.Title,
	}).Info("Blog created successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

// AddComment endpoint za dodavanje komentara
func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	var req dto.AddCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	authorID := r.Context().Value("userID").(uint)

	vars := mux.Vars(r)
	idParam := vars["id"]
	blogID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	comment, err := h.Service.AddComment(r.Context(), blogID, req, authorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// ToggleLike endpoint za lajkovanje/unlajkovanje
func (h *Handler) ToggleLike(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	vars := mux.Vars(r)
	idParam := vars["id"]
	log.WithFields(log.Fields{
		"endpoint": "/api/v1/blogs/{id}/like",
		"method":   "POST",
		"blogID":   idParam,
	}).Info("Toggle like request received")

	blogID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": "/api/v1/blogs/{id}/like",
			"blogID":   idParam,
			"error":    err.Error(),
		}).Error("Invalid blog ID format")
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	message, err := h.Service.ToggleLike(r.Context(), blogID, userID)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": "/api/v1/blogs/{id}/like",
			"blogID":   blogID.Hex(),
			"userID":   userID,
			"error":    err.Error(),
		}).Error("Failed to toggle like")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"endpoint": "/api/v1/blogs/{id}/like",
		"blogID":   blogID.Hex(),
		"userID":   userID,
		"action":   message,
	}).Info("Like toggled successfully")
	resp := map[string]string{"message": message}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

/* GetAllBlogs endpoint za dobijanje svih blogova*/

// GetAllBlogs endpoint za dobijanje svih blogova (sada je ovo feed)
func (h *Handler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"endpoint": "/api/v1/blogs",
		"method":   "GET",
	}).Info("Get all blogs request received")
    // Izvuci ID ulogovanog korisnika iz konteksta
    userID, ok := r.Context().Value("userID").(uint)
    if !ok {
		log.Error("User ID not found in context")
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

    // Pozovi novu logiku servisa
    blogs, err := h.Service.GetFeedForUser(r.Context(), userID)
    if err != nil {
		log.WithFields(log.Fields{
			"endpoint": "/api/v1/blogs",
			"userID":   userID,
			"error":    err.Error(),
		}).Error("Failed to fetch blogs")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	log.WithFields(log.Fields{
		"endpoint":   "/api/v1/blogs",
		"userID":     userID,
		"blogsCount": len(blogs),
	}).Info("Blogs fetched successfully")

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(blogs)
}

// GetBlogByID endpoint za dobijanje bloga po ID-ju
func (h *Handler) GetBlogByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]

	log.WithFields(log.Fields{
		"endpoint": "/api/v1/blogs/{id}",
		"method":   "GET",
		"blogID":   idParam,
	}).Info("Get blog by ID request received")

	blogID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": "/api/v1/blogs/{id}",
			"blogID":   idParam,
			"error":    err.Error(),
		}).Error("Invalid blog ID format")

		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	blog, err := h.Service.GetBlogByID(r.Context(), blogID)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": "/api/v1/blogs/{id}",
			"blogID":   blogID.Hex(),
			"error":    err.Error(),
		}).Error("Failed to fetch blog")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if blog == nil {
		log.WithFields(log.Fields{
			"endpoint": "/api/v1/blogs/{id}",
			"blogID":   blogID.Hex(),
		}).Warn("Blog not found")
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	log.WithFields(log.Fields{
		"endpoint": "/api/v1/blogs/{id}",
		"blogID":   blogID.Hex(),
		"title":    blog.Title,
	}).Info("Blog fetched successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}
