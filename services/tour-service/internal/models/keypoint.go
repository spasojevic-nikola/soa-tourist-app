package models

import (
	"time"
)

// KeyPoint struct represents a key point in a tour
type KeyPoint struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TourID      uint      `json:"tourId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Image       string    `json:"image"` // URL or path to image
	Order       int       `json:"order"` // Order in the tour sequence
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (KeyPoint) TableName() string { return "key_points" }