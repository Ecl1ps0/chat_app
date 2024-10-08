package auth

import (
	"ChatApp/internal/auth/models"
	"ChatApp/util"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type AuthHandler struct {
	usecase Usecase
}

func NewAuthHandler(usecase Usecase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var userData models.UserCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:        primitive.NewObjectID(),
		Username:  userData.Username,
		Password:  userData.Password,
		CreatedAt: time.Now().Unix(),
	}

	if err := h.usecase.SignUp(context.TODO(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, http.StatusOK, nil)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var userData models.UserCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := h.usecase.SignIn(context.TODO(), userData.Username, userData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authCookie := http.Cookie{
		Name:  "auth_cookie",
		Value: token,
		Path:  "/",
	}

	http.SetCookie(w, &authCookie)

	util.JSONResponse(w, http.StatusOK, map[string]string{"access_token": token})
}
