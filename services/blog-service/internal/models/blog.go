package models

import (
    "github.com/golang-jwt/jwt/v5" 
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Blog predstavlja strukturu jednog blog posta.
type Blog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	HTMLContent string 			 `bson:"htmlContent" json:"htmlContent"`
	AuthorID  uint               `bson:"authorId" json:"authorId"` // ISPRAVKA: Sada je uint
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	Images    []string           `bson:"images,omitempty" json:"images,omitempty"`
	Comments  []Comment          `bson:"comments" json:"comments"`
	Likes     []uint             `bson:"likes" json:"likes"`       // ISPRAVKA: Niz uint-ova
}

// Claims model - KLJUČNA ISPRAVKA
type Claims struct {
	UserID   uint   `json:"id"` // ISPRAVKA: Mora biti uint da se poklapa sa auth-servisom
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}