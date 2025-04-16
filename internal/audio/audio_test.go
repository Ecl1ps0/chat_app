package audio

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUsecase struct {
	mock.Mock
}

func (m *mockUsecase) CreateAudio(ctx context.Context, file multipart.File) (primitive.ObjectID, error) {
	args := m.Called(ctx, file)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *mockUsecase) GetAudio(ctx context.Context, audioId string) ([]byte, error) {
	args := m.Called(ctx, audioId)
	return args.Get(0).([]byte), args.Error(1)
}

func TestUploadAudio_Success(t *testing.T) {
	// Mock the Usecase
	usecase := new(mockUsecase)
	mockAudioID := primitive.NewObjectID() // Create a mock ObjectID
	usecase.On("CreateAudio", mock.Anything, mock.Anything).Return(mockAudioID, nil)

	// Create a handler with the mocked usecase
	handler := NewAudioHandler(usecase)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	formFile, _ := writer.CreateFormFile("audio", "test.mp3")
	formFile.Write([]byte("fake audio content"))

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a mock response recorder
	rec := httptest.NewRecorder()

	// Call the handler's UploadAudio method
	handler.UploadAudio(rec, req)

	// Assert status code and response body
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"audio_id": "`+mockAudioID.Hex()+`"}`, rec.Body.String())

	// Verify the CreateAudio method was called
	usecase.AssertExpectations(t)
}

func TestUploadAudio_MethodNotAllowed(t *testing.T) {
	usecase := new(mockUsecase)
	handler := NewAudioHandler(usecase)

	req := httptest.NewRequest(http.MethodGet, "/upload", nil)
	rec := httptest.NewRecorder()

	handler.UploadAudio(rec, req)

	// Assert method not allowed error
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestGetAudio_Success(t *testing.T) {
	// Mock the Usecase
	usecase := new(mockUsecase)
	usecase.On("GetAudio", mock.Anything, "mockAudioId").Return([]byte{0x00, 0x01, 0xFF, 0xD8, 0xFF, 0xE0}, nil)

	// Create a handler with the mocked usecase
	handler := NewAudioHandler(usecase)

	// Create a mock HTTP request
	req := httptest.NewRequest(http.MethodGet, "/audio?audio_id=mockAudioId", nil)
	rec := httptest.NewRecorder()

	// Call the handler's GetAudio method
	handler.GetAudio(rec, req)

	// Assert status code
	assert.Equal(t, http.StatusOK, rec.Code)
	// We can't directly assert the content of the file, but we can check the content-type
	assert.Equal(t, "application/octet-stream", rec.Header().Get("Content-Type"))

	// Verify the GetAudio method was called
	usecase.AssertExpectations(t)
}

func TestGetAudio_MethodNotAllowed(t *testing.T) {
	usecase := new(mockUsecase)
	handler := NewAudioHandler(usecase)

	req := httptest.NewRequest(http.MethodPost, "/audio", nil)
	rec := httptest.NewRecorder()

	handler.GetAudio(rec, req)

	// Assert method not allowed error
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestGetAudio_BadRequest(t *testing.T) {
	usecase := new(mockUsecase)
	handler := NewAudioHandler(usecase)

	req := httptest.NewRequest(http.MethodGet, "/audio", nil) // Missing audio_id
	rec := httptest.NewRecorder()

	handler.GetAudio(rec, req)

	// Assert bad request error
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "fail to get audio id\n", rec.Body.String())
}
