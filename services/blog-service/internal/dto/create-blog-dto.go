type CreateBlogDTO struct {
    Title     string   `json:"title"`
    Content   string   `json:"content"`
    Images    []string `json:"images"`
    CreatedAt string   `json:"createdAt"` // string dolazi iz frontend-a
}
