package usecase

import (
	"ChatApp/internal/message"
	models2 "ChatApp/internal/message/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MessageUsecase struct {
	repo message.Repository
}

func NewMessageUsecase(repo message.Repository) *MessageUsecase {
	return &MessageUsecase{repo: repo}
}

func (u *MessageUsecase) SaveMessage(ctx context.Context, message models2.Message) error {
	return u.repo.CreateMessage(ctx, message)
}

func (u *MessageUsecase) GetMessageByID(ctx context.Context, messageId primitive.ObjectID) (models2.Message, error) {
	return u.repo.GetMessageByID(ctx, messageId)
}

func (u *MessageUsecase) UpdateMessage(ctx context.Context, message models2.MessageDTO) error {
	return u.repo.UpdateMessage(ctx, message, time.Now().Unix())
}
