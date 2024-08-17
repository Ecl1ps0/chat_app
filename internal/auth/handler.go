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

type Handler struct {
	usecase Usecase
}

func NewAuthHandler(usecase Usecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var userData models.UserCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	user := models.User{
		ID:        primitive.NewObjectID(),
		Username:  userData.Username,
		Password:  userData.Password,
		CreatedAt: time.Now().Unix(),
	}

	userDTO, err := h.usecase.SignUp(context.TODO(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, http.StatusOK, userDTO)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var userData models.UserCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	token, err := h.usecase.SignIn(context.TODO(), userData.Username, userData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, http.StatusOK, map[string]string{"accessToken": token})
}
