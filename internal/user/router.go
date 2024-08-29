package user

import (
	"ChatApp/internal/auth"
	"net/http"
)

func (h *UserHandler) UserRouterInit(router *http.ServeMux, authHandler *auth.AuthHandler) {
	router.HandleFunc("/available-users", authHandler.AuthMiddleware(h.GetAllAvailableUsers))
}
