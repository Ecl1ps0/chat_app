package usecase

import (
	"ChatApp/internal/auth/models"
	"ChatApp/internal/user"
	"context"
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
