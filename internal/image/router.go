package image

import (
	"net/http"
)

func (h *ImageHandler) ImageRouterInit(router *http.ServeMux) {
	router.HandleFunc("/api/image", h.GetImage)
}
