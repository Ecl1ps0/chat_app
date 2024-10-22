package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ImageCode []byte             `json:"image_code" bson:"image_code"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
}
