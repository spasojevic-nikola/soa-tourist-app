package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// koristi se za dodavanje stavke u korpu.
type AddItemRequest struct {
	TourID primitive.ObjectID `json:"tourId"`
	Name   string             `json:"name"`
	Price  float64            `json:"price"`
}

//	koristi se kao odgovor nakon Checkout-a
type TourPurchaseResponse struct {
	Tokens []primitive.ObjectID `json:"purchaseTokens"`
	Message string              `json:"message"`
}