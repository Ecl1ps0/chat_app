package user

import (
	"ChatApp/util"
	"context"
	"net/http"
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

	usersDTO, err := h.usecase.GetAllUsersDTO(context.TODO())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, http.StatusOK, usersDTO)
}
