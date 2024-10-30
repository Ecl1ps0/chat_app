package user

import (
	"ChatApp/internal/auth"
	"net/http"
)

func (h *UserHandler) UserRouterInit(router *http.ServeMux, authHandler *auth.AuthHandler) {
	router.HandleFunc("/api/user/available-users", authHandler.AuthMiddleware(h.GetAllAvailableUsers))
	router.HandleFunc("/api/user", authHandler.AuthMiddleware(h.GetUser))
	router.HandleFunc("/api/update/user", authHandler.AuthMiddleware(h.UpdateUser))
}
