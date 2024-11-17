package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MessageDTO struct {
	ID       primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	SenderID primitive.ObjectID   `json:"sender_id" bson:"_id"`
	Message  string               `json:"message,omitempty"`
	Images   []primitive.ObjectID `json:"images,omitempty"`
	Audio    primitive.ObjectID   `json:"audio,omitempty"`
	IsUpdate bool                 `json:"is_update,omitempty"`
}
