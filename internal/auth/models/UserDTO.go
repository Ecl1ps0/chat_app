package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserDTO struct {
	ID             primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Username       string              `json:"username" bson:"username"`
	ProfilePicture *primitive.ObjectID `json:"profile_picture" bson:"profile_picture"`
	Bio            *string             `json:"bio" bson:"bio"`
	Email          *string             `json:"email" bson:"email"`
}
