package user

import (
	"ChatApp/util"
	"context"
	"net/http"
	"strings"
)

type UserHandler struct {
	usecase Usecase
}

func NewUserHandler(usecase Usecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) GetAllAvailableUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	tokenParts := strings.Split(token, " ")
	user, err := util.ParseToken(tokenParts[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "userId", user.ID)

	usersDTO, err := h.usecase.GetAllUsersDTO(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, http.StatusOK, usersDTO)
}
