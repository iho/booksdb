package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookStatusType string

const (
	CheckedOut BookStatusType = "CheckedOut"
	CheckedIn  BookStatusType = "CheckedIn"
)

type Book struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title,omitempty"`
	Author    string             `json:"author" bson:"author,omitempty"`
	Publisher string             `json:"publisher" bson:"publisher,omitempty"`
	Rating    int                `json:"rating" bson:"rating,omitempty"`
	Status    BookStatusType     `json:"status" bson:"status,omitempty"`
	CreatedAt time.Time          `json:"-" bson:"created_at"`
	UpdatedAt time.Time          `json:"-" bson:"updated_at"`
}
