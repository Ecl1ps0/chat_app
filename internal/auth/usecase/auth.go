package usecase

import (
	"ChatApp/internal/auth/models"
	"ChatApp/internal/auth/repository"
	"ChatApp/util"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	User models.UserDTO `json:"user"`
}

type AuthUsecase struct {
	repo repository.AuthRepository
}

func NewAuthUsecase(repo repository.AuthRepository) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

func (u *AuthUsecase) SignUp(ctx context.Context, user models.User) (models.UserDTO, error) {
	if _, err := u.repo.GetUser(ctx, user.Username); err == nil {
		return models.UserDTO{}, errors.New("such user already exist")
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
		ID:       candidate.ID,
		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		User: userDTO,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
		},
	})

	return token.SignedString([]byte(os.Getenv("SIGN_KEY")))
}

func (u *AuthUsecase) ParseToken(ctx context.Context, accessToken string) (models.UserDTO, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("SIGN_KEY")), nil
	})
	if err != nil {
		return models.UserDTO{}, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return models.UserDTO{}, errors.New("invalid access token")
	}

	return claims.User, nil
}
