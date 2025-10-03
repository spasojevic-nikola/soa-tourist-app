package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

// Enum for tour status
type TourStatus string

const (
	Draft     TourStatus = "Draft"
	Published TourStatus = "Published"
	Archived  TourStatus = "Archived"
)

// Enum for tour difficulty
type TourDifficulty string

const (
	Easy   TourDifficulty = "Easy"
	Medium TourDifficulty = "Medium"
	Hard   TourDifficulty = "Hard"
	Expert TourDifficulty = "Expert"
)

// Updated Tour struct
type Tour struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	AuthorID    uint           `json:"authorId"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Difficulty  TourDifficulty `json:"difficulty"`
	Tags        pq.StringArray `json:"tags" gorm:"type:text[]"`
	Status      TourStatus     `json:"status" gorm:"default:'Draft'"`
	Price       float64        `json:"price"`
	Distance    float64        `json:"distance" gorm:"default:0"` // Distance in kilometers
	PublishedAt *time.Time     `json:"publishedAt,omitempty"`
	ArchivedAt  *time.Time     `json:"archivedAt,omitempty"`
	IsDeleted   bool           `json:"isDeleted" gorm:"default:false"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	KeyPoints   []KeyPoint     `json:"keyPoints,omitempty" 
									gorm:"foreignKey:TourID;
									constraint:OnDelete:CASCADE"`
	Durations []TourDuration `json:"durations,omitempty" 
									gorm:"foreignKey:TourID;
									constraint:OnDelete:CASCADE"`
}

func (Tour) TableName() string { return "tours" }

type Claims struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
