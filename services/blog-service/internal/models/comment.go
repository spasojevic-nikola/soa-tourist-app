package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID  uint      `bson:"authorId" json:"authorId"`
	Text      string    `bson:"text" json:"text"`
	AuthorUsername string    `bson:"authorUsername" json:"authorUsername"`	
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}