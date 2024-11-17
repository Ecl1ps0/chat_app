package image

import (
	"net/http"
)

func (h *ImageHandler) ImageRouterInit(router *http.ServeMux, authMiddleware func(http.HandlerFunc) http.HandlerFunc) {
	router.HandleFunc("/api/image", h.GetImage)
	router.HandleFunc("/api/image/upload", authMiddleware(h.UploadImages))
}
