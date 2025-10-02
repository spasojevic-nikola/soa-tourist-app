package dto

// Novi DTO za odgovor koji sadrži SVE podatke koje frontend treba
type RecommendationDTO struct {
	UserID       uint   `json:"userId"`
	Username     string `json:"username"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	ProfileImage string `json:"profileImage"`
	Score        int    `json:"score"`
}