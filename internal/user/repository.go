package user

import (
	"ChatApp/internal/auth/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	GetAllUsersDTO(ctx context.Context) ([]models.UserDTO, error)
	GetUserById(ctx context.Context, userId primitive.ObjectID) (models.UserDTO, error)
	UpdateUser(ctx context.Context, updUser models.UserDTO) error
}
