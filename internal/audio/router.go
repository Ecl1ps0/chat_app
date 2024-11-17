package audio

import "net/http"

func (h *AudioHandler) AudioRouterInit(router *http.ServeMux, authMiddleware func(http.HandlerFunc) http.HandlerFunc) {
	router.HandleFunc("/api/audio", h.GetAudio)
	router.HandleFunc("/api/audio/upload", authMiddleware(h.UploadAudio))
}
