package models

type UserCreateDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
