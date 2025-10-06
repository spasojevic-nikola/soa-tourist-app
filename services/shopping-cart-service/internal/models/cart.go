package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//  predstavlja jednu stavku (turu) u ShoppingCart-u
type OrderItem struct {
	TourID string `bson:"tourId" json:"tourId"` // ID ture iz Tour microservice-a
	Name   string             `bson:"name" json:"name"`
	Price  float64            `bson:"price" json:"price"`
}

//  predstavlja korpu za kupovinu vezanu za jednog korisnika
type ShoppingCart struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID  uint               `bson:"userId" json:"userId"` 
	Items   []OrderItem        `bson:"items" json:"items"`
	Total   float64            `bson:"total" json:"total"`   // Ukupna cena svih stavki
	Updated time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// dodeljuje se nakon uspešne kupovine za svaku stavku
type TourPurchaseToken struct {
	ID           primitive.ObjectID 			`bson:"_id,omitempty" json:"id"`
	UserID       uint               `bson:"userId" json:"userId"`
	TourID       string `bson:"tourId" json:"tourId"`
	PurchaseTime time.Time          `bson:"purchaseTime" json:"purchaseTime"`
}

// Claims model za JWT (minimalno)
type Claims struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}