package repository

import (
	"ChatApp/internal/chat/models"
	models2 "ChatApp/internal/message/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ChatRepository struct {
	db *mongo.Collection
}

func NewChatRepository(db *mongo.Database, collection string) *ChatRepository {
	return &ChatRepository{db: db.Collection(collection)}
}

func (r *ChatRepository) CreateOrGetChat(ctx context.Context, usersIds []primitive.ObjectID) (models.Chat, error) {
	filter := bson.M{
		"members": bson.M{
			"$all": usersIds,
		},
	}

	var chat models.Chat
	if err := r.db.FindOne(ctx, filter).Decode(&chat); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			chat = models.Chat{
				ID:        primitive.NewObjectID(),
				Members:   usersIds,
				Messages:  make([]models2.Message, 0),
				CreatedAt: time.Now().Unix(),
			}

			if _, err = r.db.InsertOne(ctx, &chat); err != nil {
				return models.Chat{}, err
			}

			return chat, nil
		}

		return models.Chat{}, err
	}

	return chat, nil
}

func (r *ChatRepository) SaveMessageToChat(ctx context.Context, message models2.Message, chatId primitive.ObjectID) error {
	filter := bson.M{"_id": chatId}
	update := bson.M{"$push": bson.M{"messages": message}}

	if result := r.db.FindOneAndUpdate(ctx, filter, update); result.Err() != nil {
		return result.Err()
	}

	return nil
}
