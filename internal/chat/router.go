package chat

import (
	"net/http"
)

func (h *ChatHandler) ChatRouterInit(router *http.ServeMux) {
	router.HandleFunc("/api/chat/ws", h.StartChat)
}
