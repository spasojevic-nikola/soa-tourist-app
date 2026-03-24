package dto

type CreateTourRequest struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Difficulty  string                  `json:"difficulty"`
	Tags        []string                `json:"tags"`
	Price       float64                 `json:"price"`
	KeyPoints   []CreateKeyPointRequest `json:"keyPoints"`
}

type UpdateTourRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Difficulty  string   `json:"difficulty"`
	Tags        []string `json:"tags"`
	Price       float64  `json:"price"`
}
