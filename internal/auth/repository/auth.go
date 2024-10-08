package repository

import (
	"ChatApp/internal/auth/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	db *mongo.Collection
}

func NewAuthRepository(db *mongo.Database, collection string) *AuthRepository {
	return &AuthRepository{db: db.Collection(collection)}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user models.User) error {
	_, err := r.db.InsertOne(ctx, &user)

	return err
}

func (r *AuthRepository) GetUser(ctx context.Context, username string) (models.User, error) {
	filter := bson.D{{"username", username}}

	var user models.User
	if err := r.db.FindOne(ctx, filter).Decode(&user); err != nil {
		return models.User{}, err
	}

	return user, nil
}
