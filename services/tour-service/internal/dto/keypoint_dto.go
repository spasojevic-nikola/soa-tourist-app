package dto

// CreateKeyPointRequest DTO for creating a key point
type CreateKeyPointRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Image       string  `json:"image"`
	Order       int     `json:"order" binding:"required"`
}

// UpdateKeyPointRequest DTO for updating a key point
type UpdateKeyPointRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Image       string  `json:"image"`
	Order       int     `json:"order"`
}

// KeyPointResponse DTO for key point response
type KeyPointResponse struct {
	ID          uint    `json:"id"`
	TourID      uint    `json:"tourId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Image       string  `json:"image"`
	Order       int     `json:"order"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}