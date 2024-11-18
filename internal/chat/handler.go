package chat

import (
	"ChatApp/internal/auth/models"
	"ChatApp/internal/message"
	models2 "ChatApp/internal/message/models"
	"ChatApp/util"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"sort"
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

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Printf("fail to read message: %s", err.Error())
		return
	}

	chatId, err := primitive.ObjectIDFromHex(string(msg))
	if err != nil {
		log.Printf("fail to parse chatId: %s", err.Error())
		return
	}

	log.Printf("Chat successfully started with id: %s", chatId.Hex())

	h.mu.Lock()
	h.connections[chatId] = append(h.connections[chatId], conn)
	h.mu.Unlock()

	var messageDTO models2.MessageDTO
	for {
		_, msg, err = conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("WebSocket connection closed normally: %s", err.Error())
				break
			}

			log.Printf("fail to read a message: %s", err.Error())

			h.mu.Lock()
			conns := h.connections[chatId]
			for i, c := range conns {
				if c == conn {
					h.connections[chatId] = append(conns[:i], conns[i+1:]...)
					break
				}
			}
			h.mu.Unlock()
			return
		}

		if err = json.Unmarshal(msg, &messageDTO); err != nil {
			log.Printf("fail to unmarshal message: %s", err.Error())
			continue
		}

		h.broadcastMessage(chatId, messageDTO)
	}
}

func (h *ChatHandler) broadcastMessage(chatId primitive.ObjectID, messageDTO models2.MessageDTO) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var msg models2.Message
	var err error
	if messageDTO.IsUpdate {
		if err = h.messageUsecase.UpdateMessage(context.TODO(), messageDTO); err != nil {
			log.Printf("fail to update message: %s", err.Error())
			return
		}

		if msg, err = h.messageUsecase.GetMessageByID(context.TODO(), messageDTO.ID); err != nil {
			log.Printf("fail to get message: %s", err.Error())
			return
		}
	} else if messageDTO.DeleteFor != nil {
		if err = h.messageUsecase.DeleteMessageForUsers(context.TODO(), messageDTO.ID, messageDTO.DeleteFor); err != nil {
			log.Printf("fail to delete message: %s", err.Error())
			return
		}

		if msg, err = h.messageUsecase.GetMessageByID(context.TODO(), messageDTO.ID); err != nil {
			log.Printf("fail to get message: %s", err.Error())
			return
		}
	} else {
		msg.ID = primitive.NewObjectID()
		msg.Message = &messageDTO.Message
		msg.UserFrom = messageDTO.SenderID
		msg.ImageIDs = &messageDTO.Images
		msg.AudioID = messageDTO.Audio
		msg.DeletedFor = make(map[primitive.ObjectID]int64)
		msg.CreatedAt = time.Now().Unix()

		if err = h.messageUsecase.SaveMessage(context.TODO(), msg); err != nil {
			log.Printf("fail to save message: %s", err.Error())
			return
		}

		if err = h.chatUsecase.SaveMessageToChat(context.TODO(), msg.ID, chatId); err != nil {
			log.Printf("fail to save message to chat: %s", err.Error())
			return
		}
	}

	conns, ok := h.connections[chatId]
	if !ok {
		log.Printf("the chat with id: %s does not exist", chatId.Hex())
		return
	}

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		log.Printf("fail to marshal message: %s", err.Error())
		return
	}

	for _, conn := range conns {
		if err = conn.WriteMessage(websocket.TextMessage, msgJSON); err != nil {
			log.Printf("fail to send message: %s", err.Error())
			return
		}
	}
}

func (h *ChatHandler) ChatInit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	user1ID := r.URL.Query().Get("user1_id")
	if user1ID == "" {
		http.Error(w, "fail to get users id", http.StatusBadRequest)
		return
	}

	user2ID := r.URL.Query().Get("user2_id")
	if user2ID == "" {
		http.Error(w, "fail to get users id", http.StatusBadRequest)
		return
	}

	user1, err := primitive.ObjectIDFromHex(user1ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user2, err := primitive.ObjectIDFromHex(user2ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chat, messages, err := h.chatUsecase.CreateOrGetChat(context.TODO(), []primitive.ObjectID{user1, user2})
	if err != nil {
		log.Printf("fail to get or create chat: %s", err.Error())
		return
	}

	currentUser := r.Context().Value("currentUser").(*models.UserDTO)
	filteredMessages := make([]models2.Message, 0)
	for _, msg := range messages {
		if _, ok := msg.DeletedFor[currentUser.ID]; ok {
			continue
		}

		filteredMessages = append(filteredMessages, msg)
	}

	sort.Slice(filteredMessages, func(i, j int) bool {
		return messages[i].CreatedAt < messages[j].CreatedAt
	})

	util.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"chat_id":       chat.ID.Hex(),
		"chat_messages": filteredMessages,
	})
}
