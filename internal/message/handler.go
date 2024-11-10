package message

import (
	"ChatApp/util"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type MessageHandler struct {
	usecase Usecase
}

func NewMessageHandler(usecase Usecase) *MessageHandler {
	return &MessageHandler{usecase: usecase}
}

func (h *MessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	messageId := r.URL.Query().Get("messageId")
	if messageId == "" {
		http.Error(w, "fail to get message id", http.StatusBadRequest)
		return
	}

	hex, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message, err := h.usecase.GetMessageByID(context.TODO(), hex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, http.StatusOK, message)
}
