package models

import (
	"time"
	"github.com/lib/pq"
	"github.com/golang-jwt/jwt/v5"
)

type TourStatus string

const (
	Draft     TourStatus = "Draft"
	Published TourStatus = "Published"
	Archived  TourStatus = "Archived"
)

type Tour struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	AuthorID    uint           `json:"authorId"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Difficulty  string         `json:"difficulty"`
	Tags        pq.StringArray `json:"tags" gorm:"type:text[]"`
	Status      TourStatus     `json:"status" gorm:"default:'Draft'"`
	Price       float64        `json:"price"`
	IsDeleted   bool           `json:"isDeleted" gorm:"default:false"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

func (Tour) TableName() string { return "tours" }

type Claims struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}