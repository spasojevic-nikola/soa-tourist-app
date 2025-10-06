package models

import "time"

// TransportType enum for different transportation methods
type TransportType string

const (
	Walking TransportType = "walking"
	Bicycle TransportType = "bicycle"
	Car     TransportType = "car"
)

// TourDuration represents time needed to complete a tour by transport type
type TourDuration struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	TourID        uint          `json:"tourId"`
	TransportType TransportType `json:"transportType"`
	DurationMin   int           `json:"durationMin"` // Duration in minutes
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}

func (TourDuration) TableName() string { return "tour_durations" }
