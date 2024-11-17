package util

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ImageResponse(w http.ResponseWriter, statusCode int, imageData []byte, imageType string) {
	w.Header().Set("Content-Type", imageType)
	w.WriteHeader(statusCode)

	if _, err := w.Write(imageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
