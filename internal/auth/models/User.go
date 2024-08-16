package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
	DeletedAt int64              `json:"deleted_at" bson:"deleted_at"`
}
