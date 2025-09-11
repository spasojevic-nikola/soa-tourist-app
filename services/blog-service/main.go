package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Blog model
type Blog struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Date        time.Time          `json:"date" bson:"date"`
	Images      []string           `json:"images" bson:"images"`
	Likes       int                `json:"likes" bson:"likes"`
	Comments    []Comment          `json:"comments" bson:"comments"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Comment model
type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BlogID    primitive.ObjectID `json:"blog_id" bson:"blog_id"`
	Author    string             `json:"author" bson:"author"`
	Content   string             `json:"content" bson:"content"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// Database connection
var collection *mongo.Collection

func initDB() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	db := client.Database("blog")
	collection = db.Collection("blogs")
}

func main() {
	initDB()

	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Blog routes
	api.HandleFunc("/blogs", createBlog).Methods("POST")
	api.HandleFunc("/blogs", getAllBlogs).Methods("GET")
	api.HandleFunc("/blogs/{id}", getBlog).Methods("GET")

	// Comment routes
	api.HandleFunc("/blogs/{id}/comments", addComment).Methods("POST")
	api.HandleFunc("/blogs/{id}/comments", getComments).Methods("GET")

	// Like routes
	api.HandleFunc("/blogs/{id}/like", likeBlog).Methods("POST")
	api.HandleFunc("/blogs/{id}/unlike", unlikeBlog).Methods("POST")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	fmt.Println("Blog service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

// Create a new blog
func createBlog(w http.ResponseWriter, r *http.Request) {
	var blog Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	blog.Likes = 0
	blog.Comments = []Comment{}

	result, err := collection.InsertOne(context.TODO(), blog)
	if err != nil {
		http.Error(w, "Failed to create blog", http.StatusInternalServerError)
		return
	}

	blog.ID = result.InsertedID.(primitive.ObjectID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
}

// Get all blogs
func getAllBlogs(w http.ResponseWriter, r *http.Request) {
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch blogs", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var blogs []Blog
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		http.Error(w, "Failed to decode blogs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

// Get a specific blog
func getBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	var blog Blog
	err = collection.FindOne(context.TODO(), bson.M{"_id": blogID}).Decode(&blog)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

// Add a comment to a blog
func addComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	comment.BlogID = blogID
	comment.CreatedAt = time.Now()

	// Add comment to the blog
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": blogID},
		bson.M{"$push": bson.M{"comments": comment}},
	)
	if err != nil {
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// Get comments for a blog
func getComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	var blog Blog
	err = collection.FindOne(context.TODO(), bson.M{"_id": blogID}).Decode(&blog)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog.Comments)
}

// Like a blog
func likeBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": blogID},
		bson.M{"$inc": bson.M{"likes": 1}},
	)
	if err != nil {
		http.Error(w, "Failed to like blog", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Blog liked successfully"})
}

// Unlike a blog
func unlikeBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": blogID},
		bson.M{"$inc": bson.M{"likes": -1}},
	)
	if err != nil {
		http.Error(w, "Failed to unlike blog", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Blog unliked successfully"})
}
