package auth

import (
	"ChatApp/internal/auth/models"
	"context"
)

type Usecase interface {
	SignUp(ctx context.Context, user models.User) error
	SignIn(ctx context.Context, username, password string) (string, error)
}
