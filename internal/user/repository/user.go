package repository

import (
	"ChatApp/internal/auth/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collection string) *UserRepository {
	return &UserRepository{db: db.Collection(collection)}
}

func (r *UserRepository) GetAllUsersDTO(ctx context.Context) ([]models.UserDTO, error) {
	excludeId := ctx.Value("userId")
	filter := bson.M{"_id": bson.M{"$ne": excludeId}}

	opts := options.Find().SetProjection(bson.M{"_id": 1, "username": 1, "profile_picture": 1, "bio": 1, "email": 1})

	cursor, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var users []models.UserDTO
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, userId primitive.ObjectID) (models.UserDTO, error) {
	filter := bson.M{"_id": userId}

	opts := options.FindOne().SetProjection(bson.M{"_id": 1, "username": 1, "profile_picture": 1, "bio": 1, "email": 1})

	var user models.UserDTO
	if err := r.db.FindOne(ctx, filter, opts).Decode(&user); err != nil {
		return models.UserDTO{}, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, updUser models.UserDTO) error {
	filter := bson.M{"_id": updUser.ID}
	upd := bson.M{"$set": bson.M{"profile_picture": updUser.ProfilePicture, "bio": updUser.Bio, "email": updUser.Email}}

	_, err := r.db.UpdateOne(ctx, filter, upd)
	return err
}
