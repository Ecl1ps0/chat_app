package auth

import (
	"ChatApp/internal/auth/models"
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.User) (models.UserDTO, error)
	GetUser(ctx context.Context, username string) (models.User, error)
}
