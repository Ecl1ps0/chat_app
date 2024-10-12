package util

import (
	"ChatApp/internal/auth/models"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

const salt string = "hjqrhjqw124617ajfhajs"

type TokenClaims struct {
	jwt.RegisteredClaims
	User models.UserDTO `json:"user"`
}

func GenerateToken(user models.UserDTO) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
		},
	})

	return token.SignedString([]byte(os.Getenv("SIGN_KEY")))
}

func ParseToken(accessToken string) (models.UserDTO, error) {
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

func GenerateHash(line string) string {
	hash := sha1.New()
	hash.Write([]byte(line))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
