package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Audio struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Content   []byte             `json:"content" bson:"content"`
	CreatedAt int64              `json:"createdAt" bson:"createdAt"`
}
