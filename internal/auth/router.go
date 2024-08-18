package auth

import "net/http"

func (h *AuthHandler) AuthRouterInit(router *http.ServeMux) {
	router.HandleFunc("/sign-up", h.SignUp)
	router.HandleFunc("/sign-in", h.SignIn)
}
