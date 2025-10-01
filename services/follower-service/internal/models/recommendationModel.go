package models 

type RecommendationModel struct {
	UserID uint `json:"userId"` 
	Score  int  `json:"score"`
}