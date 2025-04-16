package image

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockImageUsecase struct {
	mock.Mock
}

func (m *mockImageUsecase) CreateImages(ctx context.Context, files []multipart.File) ([]primitive.ObjectID, error) {
	args := m.Called(ctx, files)
	return args.Get(0).([]primitive.ObjectID), args.Error(1)
}

func (m *mockImageUsecase) CreateImage(ctx context.Context, file multipart.File) (primitive.ObjectID, error) {
	args := m.Called(ctx, file)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *mockImageUsecase) GetImage(ctx context.Context, id string) ([]byte, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]byte), args.Error(1)
}

func TestUploadImages_Success(t *testing.T) {
	mockUC := new(mockImageUsecase)
	handler := NewImageHandler(mockUC)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("images", "test.jpg")
	part.Write([]byte("fake image data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	mockIDs := []primitive.ObjectID{
		primitive.NewObjectIDFromTimestamp(time.Unix(0, 0)),
		primitive.NewObjectIDFromTimestamp(time.Unix(0, 1)),
	}
	mockUC.On("CreateImages", mock.Anything, mock.Anything).Return(mockIDs, nil)

	handler.UploadImages(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	expectedJSON, _ := json.Marshal(mockIDs)
	assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	mockUC.AssertExpectations(t)
}

func TestUploadImages_NoFiles(t *testing.T) {
	mockUC := new(mockImageUsecase)
	handler := NewImageHandler(mockUC)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	handler.UploadImages(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Contains(t, rec.Body.String(), "No images uploaded")
}

func TestUploadImages_MethodNotAllowed(t *testing.T) {
	mockUC := new(mockImageUsecase)
	handler := NewImageHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/upload", nil)
	rec := httptest.NewRecorder()

	handler.UploadImages(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestGetImage_MissingID(t *testing.T) {
	mockUC := new(mockImageUsecase)
	handler := NewImageHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/image", nil)
	rec := httptest.NewRecorder()

	handler.GetImage(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Fail to get image id\n", rec.Body.String())
}

func TestGetImage_MethodNotAllowed(t *testing.T) {
	mockUC := new(mockImageUsecase)
	handler := NewImageHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/image", nil)
	rec := httptest.NewRecorder()

	handler.GetImage(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}
