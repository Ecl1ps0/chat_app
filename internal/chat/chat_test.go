package chat

import (
	"ChatApp/internal/auth/models"
	models2 "ChatApp/internal/chat/models"
	models3 "ChatApp/internal/message/models"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockChatUsecase struct {
	mock.Mock
}

func (m *mockChatUsecase) CreateOrGetChat(ctx context.Context, userIDs []primitive.ObjectID) (models2.Chat, []models3.Message, error) {
	args := m.Called(ctx, userIDs)
	return args.Get(0).(models2.Chat), args.Get(1).([]models3.Message), args.Error(2)
}

func (m *mockChatUsecase) SaveMessageToChat(ctx context.Context, messageID primitive.ObjectID, chatID primitive.ObjectID) error {
	args := m.Called(ctx, messageID, chatID)
	return args.Error(0)
}

type mockMessageUsecase struct {
	mock.Mock
}

func (m *mockMessageUsecase) SaveMessage(ctx context.Context, msg models3.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *mockMessageUsecase) GetMessageByID(ctx context.Context, id primitive.ObjectID) (models3.Message, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models3.Message), args.Error(1)
}

func (m *mockMessageUsecase) UpdateMessage(ctx context.Context, dto models3.MessageDTO) error {
	args := m.Called(ctx, dto)
	return args.Error(0)
}

func (m *mockMessageUsecase) DeleteMessageForUsers(ctx context.Context, id primitive.ObjectID, users []primitive.ObjectID) error {
	args := m.Called(ctx, id, users)
	return args.Error(0)
}

func TestChatInit_Success(t *testing.T) {
	mockChat := new(mockChatUsecase)
	mockMessage := new(mockMessageUsecase)

	user1ID := primitive.NewObjectID()
	user2ID := primitive.NewObjectID()
	chatID := primitive.NewObjectID()

	mockedChat := models2.Chat{ID: chatID}
	msg := "Hello"
	mockedMessages := []models3.Message{
		{
			ID:         primitive.NewObjectID(),
			Message:    &msg,
			UserFrom:   user1ID,
			CreatedAt:  1,
			DeletedFor: make(map[primitive.ObjectID]int64),
		},
	}

	mockChat.On("CreateOrGetChat", mock.Anything, mock.Anything).Return(mockedChat, mockedMessages, nil)

	h := NewChatHandler(mockChat, mockMessage)

	req := httptest.NewRequest(http.MethodGet, "/chat/init?user1_id="+user1ID.Hex()+"&user2_id="+user2ID.Hex(), nil)
	req = req.WithContext(context.WithValue(req.Context(), "currentUser", &models.UserDTO{ID: user1ID}))

	rec := httptest.NewRecorder()

	h.ChatInit(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string]interface{}
	err := json.NewDecoder(rec.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, chatID.Hex(), body["chat_id"])
	assert.Len(t, body["chat_messages"], 1)
}

func TestChatInit_BadRequest(t *testing.T) {
	h := NewChatHandler(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/chat/init", nil)
	rec := httptest.NewRecorder()

	h.ChatInit(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/chat/init?user1_id=abc", nil)
	rec = httptest.NewRecorder()
	// adding user1_id only, missing user2_id
	h.ChatInit(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestChatInit_InternalServerError(t *testing.T) {
	mockChat := new(mockChatUsecase)
	h := NewChatHandler(mockChat, nil)

	user1ID := primitive.NewObjectID().Hex()
	user2ID := "invalidHex"

	req := httptest.NewRequest(http.MethodGet, "/chat/init?user1_id="+user1ID+"&user2_id="+user2ID, nil)
	rec := httptest.NewRecorder()

	h.ChatInit(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
