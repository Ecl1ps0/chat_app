package usecase

import (
	"ChatApp/internal/chat/models"
	"ChatApp/internal/chat/repository"
	models2 "ChatApp/internal/message/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatUsecase struct {
	repo repository.ChatRepository
}

func NewChatUsecase(repo repository.ChatRepository) *ChatUsecase {
	return &ChatUsecase{repo: repo}
}

func (u *ChatUsecase) CreateOrGetChat(ctx context.Context, usersIds []primitive.ObjectID) (models.Chat, error) {
	return u.repo.CreateOrGetChat(ctx, usersIds)
}

func (u *ChatUsecase) SaveMessageToChat(ctx context.Context, message models2.Message, chatId primitive.ObjectID) error {
	return u.repo.SaveMessageToChat(ctx, message, chatId)
}