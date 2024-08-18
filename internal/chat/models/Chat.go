package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Members   []primitive.ObjectID `json:"members" bson:"members"`
	Messages  []primitive.ObjectID `json:"messages" bson:"messages"`
	CreatedAt int64                `json:"created_at" bson:"created_at"`
}
