package audio

import (
	"ChatApp/util"
	"context"
	"net/http"
)

type AudioHandler struct {
	usecase Usecase
}

func NewAudioHandler(usecase Usecase) *AudioHandler {
	return &AudioHandler{usecase}
}

func (h *AudioHandler) UploadAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audioId, err := h.usecase.CreateAudio(context.TODO(), file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"audio_id": audioId.Hex(),
	})
}

func (h *AudioHandler) GetAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	audioId := r.URL.Query().Get("audio_id")
	if audioId == "" {
		http.Error(w, "fail to get audio id", http.StatusBadRequest)
		return
	}

	audio, err := h.usecase.GetAudio(context.TODO(), audioId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	audioType := http.DetectContentType(audio)

	util.FileResponse(w, http.StatusOK, audio, audioType)
}
