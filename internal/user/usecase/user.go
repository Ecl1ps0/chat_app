package usecase

import (
	"ChatApp/internal/auth/models"
	"ChatApp/internal/user"
	"ChatApp/util"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUsecase struct {
	repo user.Repository
}

func NewUserUsecase(repo user.Repository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetAllUsersDTO(ctx context.Context) ([]models.UserDTO, error) {
	return u.repo.GetAllUsersDTO(ctx)
}

func (u *UserUsecase) GetUserById(ctx context.Context, userId string) (models.UserDTO, error) {
	hex, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return models.UserDTO{}, err
	}
	return u.repo.GetUserById(ctx, hex)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, updUser models.UserDTO) (string, error) {
	if err := u.repo.UpdateUser(ctx, updUser); err != nil {
		return "", err
	}

	return util.GenerateToken(updUser)
}
