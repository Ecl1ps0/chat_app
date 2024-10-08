package message

import (
	models2 "ChatApp/internal/message/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Usecase interface {
	SaveMessage(ctx context.Context, message models2.Message) error
	GetMessageByID(ctx context.Context, messageId primitive.ObjectID) (models2.Message, error)
}
