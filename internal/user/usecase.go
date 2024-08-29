package user

import (
	"ChatApp/internal/auth/models"
	"context"
)

type Usecase interface {
	GetAllUsersDTO(ctx context.Context) ([]models.UserDTO, error)
}
