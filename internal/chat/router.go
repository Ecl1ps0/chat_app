package chat

import (
	"net/http"
)

func (h *ChatHandler) ChatRouterInit(router *http.ServeMux) {
	router.HandleFunc("/chat-ws", h.StartChat)
}
