package user

import (
	"ChatApp/internal/auth/models"
	"context"
)

type Usecase interface {
	GetAllUsersDTO(ctx context.Context) ([]models.UserDTO, error)
	GetUserById(ctx context.Context, userId string) (models.UserDTO, error)
	UpdateUser(ctx context.Context, updUser models.UserDTO) (string, error)
}
