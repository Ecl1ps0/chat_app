package chat

import (
	"ChatApp/internal/chat/models"
	models2 "ChatApp/internal/message/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Usecase interface {
	CreateOrGetChat(ctx context.Context, usersIds []primitive.ObjectID) (models.Chat, error)
	SaveMessageToChat(ctx context.Context, message models2.Message, chatId primitive.ObjectID) error
}
