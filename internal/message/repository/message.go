package repository

import (
	models2 "ChatApp/internal/message/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepository struct {
	db *mongo.Collection
}

func NewMessageRepository(db *mongo.Database, collection string) *MessageRepository {
	return &MessageRepository{db: db.Collection(collection)}
}

func (r *MessageRepository) CreateMessage(ctx context.Context, message models2.Message) error {
	if _, err := r.db.InsertOne(ctx, &message); err != nil {
		return err
	}

	return nil
}

func (r *MessageRepository) GetMessageByID(ctx context.Context, messageId primitive.ObjectID) (models2.Message, error) {
	filter := bson.D{{"_id", messageId}}

	var message models2.Message
	if err := r.db.FindOne(ctx, filter).Decode(&message); err != nil {
		return models2.Message{}, err
	}

	return message, nil
}
