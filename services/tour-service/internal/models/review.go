package models

import (
	"time"
)

// Review struct represents a review for a tour
type Review struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	TourID          uint      `json:"tourId" gorm:"not null"`
	TouristID       uint      `json:"touristId" gorm:"not null"`                // ID korisnika koji je ostavio recenziju
	TouristUsername string    `json:"touristUsername" gorm:"type:varchar(255)"` // Username turiste
	Rating          int       `json:"rating" gorm:"not null"`                   // Ocena 1-5
	Comment         string    `json:"comment" gorm:"type:text"`
	VisitDate       time.Time `json:"visitDate"`               // Datum kada je posetio turu
	Images          string    `json:"images" gorm:"type:text"` // JSON array URLs slika
	CreatedAt       time.Time `json:"createdAt"`               // Datum kada je ostavio komentar
	UpdatedAt       time.Time `json:"updatedAt"`

	// Relacije
	Tour Tour `json:"tour,omitempty" gorm:"foreignKey:TourID"`
}

func (Review) TableName() string { return "reviews" }
