package message

import "net/http"

func (h *MessageHandler) MessageRouterInit(router *http.ServeMux, authMiddleware func(http.HandlerFunc) http.HandlerFunc) {
	router.HandleFunc("/api/message", authMiddleware(h.GetMessage))
}
