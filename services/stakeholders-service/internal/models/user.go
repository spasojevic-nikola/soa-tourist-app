package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User model predstavlja korisnika u stakeholders_users tabeli
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"unique"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	ProfileImage string    `json:"profile_image"`
	Biography    string    `json:"biography"`
	Motto        string    `json:"motto"`
	Role         string    `json:"role" gorm:"default:'tourist'"`
	IsBlocked    bool      `json:"is_blocked" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName specificira ime tabele za User model
func (User) TableName() string {
	return "stakeholders_users"
}

// Claims model predstavlja podatke koji se čuvaju unutar JWT tokena
type Claims struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}