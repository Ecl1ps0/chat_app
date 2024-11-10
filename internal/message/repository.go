package message

import (
	models2 "ChatApp/internal/message/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	CreateMessage(ctx context.Context, message models2.Message) error
	GetMessageByID(ctx context.Context, messageId primitive.ObjectID) (models2.Message, error)
	UpdateMessage(ctx context.Context, message models2.MessageDTO, updateTime int64) error
}
