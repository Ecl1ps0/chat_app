package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         primitive.ObjectID           `json:"id" bson:"_id,omitempty"`
	Message    *string                      `json:"message,omitempty" bson:"message"`
	UserFrom   primitive.ObjectID           `json:"user_from" bson:"user_from"`
	ImageIDs   *[]primitive.ObjectID        `json:"images,omitempty" bson:"images"`
	AudioID    *primitive.ObjectID          `json:"audio,omitempty" bson:"audio"`
	DeletedFor map[primitive.ObjectID]int64 `json:"deleted_for,omitempty" bson:"deleted_for,omitempty"`
	CreatedAt  int64                        `json:"created_at" bson:"created_at"`
	UpdatedAt  int64                        `json:"updated_at" bson:"updated_at"`
}
