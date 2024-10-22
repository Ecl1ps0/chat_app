package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Message   *string               `json:"message,omitempty" bson:"message"`
	UserFrom  primitive.ObjectID    `json:"user_from" bson:"user_from"`
	ImageIDs  *[]primitive.ObjectID `json:"images,omitempty" bson:"images"`
	CreatedAt int64                 `json:"created_at" bson:"created_at"`
}
