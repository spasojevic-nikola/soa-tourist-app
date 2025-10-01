package api

import (
	"encoding/json"
	"net/http"

	"blog-service/internal/dto"
	"blog-service/internal/service"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	var req dto.CreateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	authorID := r.Context().Value("userID").(uint)

	blog, err := h.Service.CreateBlog(r.Context(), req, authorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	blogID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	message, err := h.Service.ToggleLike(r.Context(), blogID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"message": message}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

/* GetAllBlogs endpoint za dobijanje svih blogova*/

// GetAllBlogs endpoint za dobijanje svih blogova (sada je ovo feed)
func (h *Handler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
    // Izvuci ID ulogovanog korisnika iz konteksta
    userID, ok := r.Context().Value("userID").(uint)
    if !ok {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

    // Pozovi novu logiku servisa
    blogs, err := h.Service.GetFeedForUser(r.Context(), userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(blogs)
}

// GetBlogByID endpoint za dobijanje bloga po ID-ju
func (h *Handler) GetBlogByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	blogID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	blog, err := h.Service.GetBlogByID(r.Context(), blogID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if blog == nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}
