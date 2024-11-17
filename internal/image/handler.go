package image

import (
	"ChatApp/util"
	"context"
	"mime/multipart"
	"net/http"
)

type ImageHandler struct {
	usecase Usecase
}

func NewImageHandler(usecase Usecase) *ImageHandler {
	return &ImageHandler{usecase: usecase}
}

func (h *ImageHandler) UploadImages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		util.JSONResponse(w, http.StatusNoContent, "No images uploaded")
		return
	}

	var images []multipart.File
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		images = append(images, file)
		file.Close()
	}

	imageIds, err := h.usecase.CreateImages(context.TODO(), images)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, http.StatusOK, imageIds)
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

	util.FileResponse(w, http.StatusOK, imageData, imageType)
}
