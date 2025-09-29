package dto

type CreateBlogRequest struct {
	Title 	string 	`json:"title" validate:"required"`
	Content string 	`json:"content" validate:"required"`
	Images 	[]string `json:"images,omitempty"`
}

type AddCommentRequest struct {
	Text string `json:"text" validate:"required"`
}
