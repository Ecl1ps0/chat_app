package usecase

import (
	"ChatApp/internal/chat"
	"ChatApp/internal/chat/models"
	models2 "ChatApp/internal/message/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatUsecase struct {
	repo chat.Repository
}

func NewChatUsecase(repo chat.Repository) *ChatUsecase {
	return &ChatUsecase{repo: repo}
}

func (u *ChatUsecase) CreateOrGetChat(ctx context.Context, usersIds []primitive.ObjectID) (models.Chat, error) {
	return u.repo.CreateOrGetChat(ctx, usersIds)
}

func (u *ChatUsecase) SaveMessageToChat(ctx context.Context, message models2.Message, chatId primitive.ObjectID) error {
	return u.repo.SaveMessageToChat(ctx, message, chatId)
}
