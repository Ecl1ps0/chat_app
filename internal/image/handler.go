package image

import (
	"context"
	"net/http"
)

type ImageHandler struct {
	usecase Usecase
}

func NewImageHandler(usecase Usecase) *ImageHandler {
	return &ImageHandler{usecase: usecase}
}

func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	imageId := r.URL.Query().Get("id")
	if imageId == "" {
		http.Error(w, "Fail to get image id", http.StatusBadRequest)
		return
	}

	imageData, err := h.usecase.GetImage(context.TODO(), imageId)
	if err != nil {
		http.Error(w, "Fail to get image", http.StatusInternalServerError)
		return
	}

	imageType := http.DetectContentType(imageData)

	w.Header().Set("Content-Type", imageType)
	if _, err := w.Write(imageData); err != nil {
		http.Error(w, "Fail to write image", http.StatusInternalServerError)
		return
	}
}
