package auth

import "net/http"

func (h *Handler) AuthRouterInit(router *http.ServeMux) {
	router.HandleFunc("/sign-up", h.SignUp)
	router.HandleFunc("/sign-in", h.SignIn)
}
