package usecase

import (
	"ChatApp/internal/auth"
	"ChatApp/internal/auth/models"
	"ChatApp/util"
	"context"
	"errors"
)

type AuthUsecase struct {
	repo auth.Repository
}

func NewAuthUsecase(repo auth.Repository) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

func (u *AuthUsecase) SignUp(ctx context.Context, user models.User) error {
	if _, err := u.repo.GetUser(ctx, user.Username); err == nil {
		return errors.New("such user already exist")
	}

	user.Password = util.GenerateHash(user.Password)

	return u.repo.CreateUser(ctx, user)
}

func (u *AuthUsecase) SignIn(ctx context.Context, username, password string) (string, error) {
	candidate, err := u.repo.GetUser(ctx, username)
	if err != nil {
		return "", err
	}

	if candidate.Password != util.GenerateHash(password) {
		return "", errors.New("wrong password")
	}

	userDTO := models.UserDTO{
		ID:             candidate.ID,
		Username:       username,
		ProfilePicture: candidate.ProfilePicture,
		Bio:            candidate.Bio,
		Email:          candidate.Email,
	}

	return util.GenerateToken(userDTO)
}
