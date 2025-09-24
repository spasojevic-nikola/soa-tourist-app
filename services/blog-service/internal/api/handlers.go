package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"blog-service/internal/models" 

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	DB *mongo.Database
}

func NewHandler(db *mongo.Database) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	// Citamo userID kao uint iz konteksta
	authorID, _ := r.Context().Value("userID").(uint)

	var blog models.Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	blog.ID = primitive.NewObjectID()
	blog.AuthorID = authorID
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	blog.Comments = []models.Comment{}
	blog.Likes = []uint{}

	_, err := h.DB.Collection("blogs").InsertOne(context.Background(), blog)
	if err != nil {
		http.Error(w, "Failed to create blog", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
}

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	blogID, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	//Citamo userID kao uint
	authorID, _ := r.Context().Value("userID").(uint)

	var commentData struct{ Text string `json:"text"` }
	json.NewDecoder(r.Body).Decode(&commentData)

	newComment := models.Comment{
		AuthorID:  authorID,
		Text:      commentData.Text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	update := bson.M{"$push": bson.M{"comments": newComment}}
	_, err := h.DB.Collection("blogs").UpdateOne(context.Background(), bson.M{"_id": blogID}, update)
	if err != nil {
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newComment)
}

func (h *Handler) ToggleLike(w http.ResponseWriter, r *http.Request) {
	blogID, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	// Citamo userID kao uint
	userID, _ := r.Context().Value("userID").(uint)

	var blog models.Blog
	h.DB.Collection("blogs").FindOne(context.Background(), bson.M{"_id": blogID}).Decode(&blog)

	alreadyLiked := false
	for _, id := range blog.Likes {
		if id == userID {
			alreadyLiked = true
			break
		}
	}
	
	var update bson.M
	message := ""
	if alreadyLiked {
		update = bson.M{"$pull": bson.M{"likes": userID}}
		message = "Blog unliked successfully"
	} else {
		update = bson.M{"$addToSet": bson.M{"likes": userID}}
		message = "Blog liked successfully"
	}

	_, err := h.DB.Collection("blogs").UpdateOne(context.Background(), bson.M{"_id": blogID}, update)
	if err != nil {
		http.Error(w, "Failed to update like status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}