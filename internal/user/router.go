package user

import (
	"net/http"
)

func (h *UserHandler) UserRouterInit(router *http.ServeMux, authMiddleware func(http.HandlerFunc) http.HandlerFunc) {
	router.HandleFunc("/api/user/available-users", authMiddleware(h.GetAllAvailableUsers))
	router.HandleFunc("/api/user", authMiddleware(h.GetUser))
	router.HandleFunc("/api/update/user", authMiddleware(h.UpdateUser))
}
