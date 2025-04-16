package message

import (
	"ChatApp/internal/message/models"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockUsecase struct {
	mock.Mock
}

func (m *MockUsecase) GetMessageByID(ctx context.Context, id primitive.ObjectID) (models.Message, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Message), args.Error(1)
}

func (m *MockUsecase) SaveMessage(ctx context.Context, msg models.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MockUsecase) UpdateMessage(ctx context.Context, msg models.MessageDTO) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MockUsecase) DeleteMessageForUsers(ctx context.Context, msgID primitive.ObjectID, users []primitive.ObjectID) error {
	args := m.Called(ctx, msgID, users)
	return args.Error(0)
}

func TestGetMessage_Success(t *testing.T) {
	mockUC := new(MockUsecase)
	handler := NewMessageHandler(mockUC)

	messageID := primitive.NewObjectID()
	expectedMessage := models.Message{
		ID:        messageID,
		Message:   strPtr("Hello, world!"),
		UserFrom:  primitive.NewObjectID(),
		CreatedAt: 1234567890,
	}

	mockUC.On("GetMessageByID", mock.Anything, messageID).Return(expectedMessage, nil)

	req := httptest.NewRequest(http.MethodGet, "/message?messageId="+messageID.Hex(), nil)
	rec := httptest.NewRecorder()

	handler.GetMessage(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Hello, world!")
	mockUC.AssertExpectations(t)
}

func TestGetMessage_InvalidID(t *testing.T) {
	mockUC := new(MockUsecase)
	handler := NewMessageHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/message?messageId=invalid_id", nil)
	rec := httptest.NewRecorder()

	handler.GetMessage(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetMessage_NotFound(t *testing.T) {
	mockUC := new(MockUsecase)
	handler := NewMessageHandler(mockUC)

	msgID := primitive.NewObjectID()
	mockUC.On("GetMessageByID", mock.Anything, msgID).Return(models.Message{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/message?messageId="+msgID.Hex(), nil)
	rec := httptest.NewRecorder()

	handler.GetMessage(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "not found")
}

func TestGetMessage_MethodNotAllowed(t *testing.T) {
	mockUC := new(MockUsecase)
	handler := NewMessageHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/message?messageId=someid", nil)
	rec := httptest.NewRecorder()

	handler.GetMessage(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func strPtr(s string) *string {
	return &s
}
