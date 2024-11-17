package chat

import (
	"net/http"
)

func (h *ChatHandler) ChatRouterInit(router *http.ServeMux, authMiddleware func(http.HandlerFunc) http.HandlerFunc) {
	router.HandleFunc("/chat/ws", h.StartChat)
	router.HandleFunc("/api/chat/init", authMiddleware(h.ChatInit))
}
