package dto

// AddDurationRequest DTO for adding tour duration
type AddDurationRequest struct {
	TransportType string `json:"transportType"` // "walking", "bicycle", "car"
	DurationMin   int    `json:"durationMin"`
}

// DurationResponse DTO for duration response
type DurationResponse struct {
	ID            uint   `json:"id"`
	TourID        uint   `json:"tourId"`
	TransportType string `json:"transportType"`
	DurationMin   int    `json:"durationMin"`
}
