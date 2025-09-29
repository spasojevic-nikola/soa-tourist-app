package models

import (
	"time"
)

type Comment struct {
	AuthorID  uint      `bson:"authorId" json:"authorId"`
	Text      string    `bson:"text" json:"text"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}