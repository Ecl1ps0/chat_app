package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Username       string              `json:"username" bson:"username"`
	Password       string              `json:"password" bson:"password"`
	ProfilePicture *primitive.ObjectID `json:"profile_picture" bson:"profile_picture"`
	Bio            *string             `json:"bio" bson:"bio"`
	Email          *string             `json:"email" bson:"email"`
	CreatedAt      int64               `json:"created_at" bson:"created_at"`
	UpdatedAt      int64               `json:"updated_at" bson:"updated_at"`
	DeletedAt      int64               `json:"deleted_at" bson:"deleted_at"`
}
