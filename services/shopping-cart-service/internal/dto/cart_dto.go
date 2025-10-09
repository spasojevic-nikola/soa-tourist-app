package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// koristi se za dodavanje stavke u korpu.
type AddItemRequest struct {
	TourID string `json:"tourId"`
}

//	koristi se kao odgovor nakon Checkout-a
type TourPurchaseResponse struct {
	Tokens []primitive.ObjectID `json:"purchaseTokens"`
	Message string              `json:"message"`
}