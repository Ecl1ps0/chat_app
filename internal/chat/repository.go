package chat

import (
	"ChatApp/internal/chat/models"
	models2 "ChatApp/internal/message/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	CreateOrGetChat(ctx context.Context, usersIds []primitive.ObjectID) (models.Chat, []models2.Message, error)
	SaveMessageToChat(ctx context.Context, messageId primitive.ObjectID, chatId primitive.ObjectID) error
}
