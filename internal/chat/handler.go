package chat

import (
	"ChatApp/internal/chat/models"
	"ChatApp/internal/message"
	models2 "ChatApp/internal/message/models"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"sync"
	"time"
)

type ChatHandler struct {
	chatUsecase    Usecase
	messageUsecase message.Usecase
	connections    map[primitive.ObjectID][]*websocket.Conn
	mu             sync.Mutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewChatHandler(chatUsecase Usecase, messageUsecase message.Usecase) *ChatHandler {
	return &ChatHandler{
		chatUsecase:    chatUsecase,
		messageUsecase: messageUsecase,
		connections:    make(map[primitive.ObjectID][]*websocket.Conn),
	}
}

func (h *ChatHandler) StartChat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("fail to upgrade to websocket: %s", err.Error())
		return
	}
	defer conn.Close()

	var usersIds models.UsersIdsDTO
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Fatalf("fail to read message: %s", err.Error())
		return
	}

	if err = json.Unmarshal(msg, &usersIds); err != nil {
		log.Fatalf("fail to unmarshal: %s", err.Error())
		return
	}

	user1ID, err := primitive.ObjectIDFromHex(usersIds.User1ID)
	if err != nil {
		log.Fatalf("invalid user1 ID: %s", err.Error())
		return
	}

	user2ID, err := primitive.ObjectIDFromHex(usersIds.User2ID)
	if err != nil {
		log.Fatalf("invalid user2 ID: %s", err.Error())
		return
	}

	chat, err := h.chatUsecase.CreateOrGetChat(context.TODO(), []primitive.ObjectID{user1ID, user2ID})
	if err != nil {
		log.Fatalf("fail to get or create chat: %s", err.Error())
		return
	}

	log.Printf("Chat successfully started with id: %s", chat.ID.Hex())

	h.mu.Lock()
	h.connections[chat.ID] = append(h.connections[chat.ID], conn)
	h.mu.Unlock()

	var messageDTO models2.MessageDTO
	for {
		_, msg, err = conn.ReadMessage()
		if err != nil {
			log.Fatalf("fail to read a message: %s", err.Error())
			return
		}

		if err = json.Unmarshal(msg, &messageDTO); err != nil {
			log.Fatalf("fail to unmarshal message: %s", err.Error())
			return
		}

		h.broadcastMessage(&chat, messageDTO.SenderID, messageDTO.Message)
	}
}

func (h *ChatHandler) broadcastMessage(chat *models.Chat, senderId, message string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sender, err := primitive.ObjectIDFromHex(senderId)
	if err != nil {
		log.Fatalf("invalid sender id: %s", err.Error())
		return
	}

	msg := models2.Message{
		ID:        primitive.NewObjectID(),
		Message:   message,
		UserFrom:  sender,
		CreatedAt: time.Now().Unix(),
	}

	if err = h.messageUsecase.SaveMessage(context.TODO(), msg); err != nil {
		log.Fatalf("fail to save message: %s", err.Error())
		return
	}

	if err = h.chatUsecase.SaveMessageToChat(context.TODO(), msg, chat.ID); err != nil {
		log.Fatalf("fail to save message to chat: %s", err.Error())
		return
	}

	conns, ok := h.connections[chat.ID]
	if !ok {
		log.Fatalf("the chat with id: %s does not exist", chat.ID.Hex())
		return
	}

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("fail to marshal message: %s", err.Error())
		return
	}

	for _, conn := range conns {
		if err = conn.WriteMessage(websocket.TextMessage, msgJSON); err != nil {
			log.Fatalf("fail to send message: %s", err.Error())
			return
		}
	}
}
