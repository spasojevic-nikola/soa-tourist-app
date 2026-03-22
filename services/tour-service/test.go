package main
import (
"encoding/json"
"fmt"
"time"
)
type CreateReviewRequest struct {
TourID    uint      json:"tourId" binding:"required"
Rating    int       json:"rating" binding:"required,min=1,max=5"
Comment   string    json:"comment"
VisitDate time.Time json:"visitDate" binding:"required"
Images    []string  json:"images"
}
func main() {
var req CreateReviewRequest
js := []byte( + "" + {"tourId": 1, "rating": 5, "comment": "abc", "visitDate": "2024-10-25T14:48:00.000Z", "images": []} + "" + )
err := json.Unmarshal(js, &req)
fmt.Println(err, req)
}
