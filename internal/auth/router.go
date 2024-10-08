package auth

import "net/http"

func (h *AuthHandler) AuthRouterInit(router *http.ServeMux) {
	router.HandleFunc("/api/auth/sign-up", h.SignUp)
	router.HandleFunc("/api/auth/sign-in", h.SignIn)
}
