package repository

import (
	"ChatApp/internal/chat/models"
	models2 "ChatApp/internal/message/models"
	"context"
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

func (r *ChatRepository) CreateOrGetChat(ctx context.Context, usersIds []primitive.ObjectID) (models.Chat, []models2.Message, error) {
	pipline := mongo.Pipeline{
		{{"$match", bson.M{
			"members": bson.M{
				"$all": usersIds,
			},
		}}},
		{{"$lookup", bson.M{
			"from":         "messages",
			"localField":   "messages",
			"foreignField": "_id",
			"as":           "message_details",
		}}},
	}

	cursor, err := r.db.Aggregate(ctx, pipline)
	if err != nil {
		return models.Chat{}, nil, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Chat            models.Chat       `bson:",inline"`
		MessagesDetails []models2.Message `bson:"message_details"`
	}
	if cursor.Next(ctx) {
		if err = cursor.Decode(&result); err != nil {
			return models.Chat{}, nil, err
		}

		return result.Chat, result.MessagesDetails, nil
	}
	if cursor.Err() != nil {
		return models.Chat{}, nil, err
	}

	chat := models.Chat{
		ID:        primitive.NewObjectID(),
		Members:   usersIds,
		Messages:  make([]primitive.ObjectID, 0),
		CreatedAt: time.Now().Unix(),
	}

	if _, err = r.db.InsertOne(ctx, &chat); err != nil {
		return models.Chat{}, nil, err
	}

	return chat, nil, nil
}

func (r *ChatRepository) SaveMessageToChat(ctx context.Context, messageId primitive.ObjectID, chatId primitive.ObjectID) error {
	filter := bson.M{"_id": chatId}
	update := bson.M{"$push": bson.M{"messages": messageId}}

	if result := r.db.FindOneAndUpdate(ctx, filter, update); result.Err() != nil {
		return result.Err()
	}

	return nil
}
