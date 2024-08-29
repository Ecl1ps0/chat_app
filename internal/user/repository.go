package user

import (
	"ChatApp/internal/auth/models"
	"context"
)

type Repository interface {
	GetAllUsersDTO(ctx context.Context) ([]models.UserDTO, error)
}
