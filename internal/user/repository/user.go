package repository

import (
	"ChatApp/internal/auth/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
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
	opts := options.Find().SetProjection(bson.M{"_id": 1, "username": 1})

	cursor, err := r.db.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	var users []models.UserDTO
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
