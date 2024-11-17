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
	_, err := r.db.InsertOne(ctx, &message)

	return err
}

func (r *MessageRepository) GetMessageByID(ctx context.Context, messageId primitive.ObjectID) (models2.Message, error) {
	filter := bson.M{"_id": messageId}

	var message models2.Message
	if err := r.db.FindOne(ctx, filter).Decode(&message); err != nil {
		return models2.Message{}, err
	}

	return message, nil
}

func (r *MessageRepository) UpdateMessage(ctx context.Context, message models2.MessageDTO, updateTime int64) error {
	filter := bson.D{{"_id", message.ID}}
	upd := bson.M{"$set": bson.M{"message": message.Message, "updated_at": updateTime}}

	_, err := r.db.UpdateOne(ctx, filter, upd)
	return err
}
