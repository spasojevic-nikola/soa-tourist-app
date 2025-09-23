package dto

type UserOverviewDto struct {
    ID          uint   `json:"id"`
    Username    string `json:"username"`
    Email       string `json:"email"`
    FirstName   string `json:"first_name"`
    LastName    string `json:"last_name"`
    ProfileImage string `json:"profile_image"`
    Role        string `json:"role"`
}
