package user

import (
	"ChatApp/internal/auth/models"
	"ChatApp/util"
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserUsecase struct {
	mock.Mock
}

type mockImageUsecase struct {
	mock.Mock
}

func (m *mockUserUsecase) GetUserById(ctx context.Context, userId string) (models.UserDTO, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(models.UserDTO), args.Error(1)
}

func (m *mockUserUsecase) GetAllUsersDTO(ctx context.Context) ([]models.UserDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.UserDTO), args.Error(1)
}

func (m *mockUserUsecase) UpdateUser(ctx context.Context, user models.UserDTO) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
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

func TestGetAllAvailableUsers_Success(t *testing.T) {
	mockUC := new(mockUserUsecase)
	mockImageUC := new(mockImageUsecase)
	handler := NewUserHandler(mockUC, mockImageUC)

	userDTO := models.UserDTO{ID: primitive.NewObjectID(), Username: "john"}

	// Generate a valid token using the GenerateToken function
	token, err := util.GenerateToken(userDTO)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}

	mockUsers := []models.UserDTO{{ID: userDTO.ID, Username: "john"}}

	// Ensure the mock returns no error
	mockUC.On("GetAllUsersDTO", mock.Anything).Return(mockUsers, nil).Once()

	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	r.Header.Set("Authorization", "Bearer "+token) // Use the valid token here
	w := httptest.NewRecorder()

	// Call the handler method
	handler.GetAllAvailableUsers(w, r)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}
}

func TestGetUser_Success(t *testing.T) {
	mockUC := new(mockUserUsecase)
	mockImageUC := new(mockImageUsecase)
	handler := NewUserHandler(mockUC, mockImageUC)

	userId := primitive.NewObjectID().Hex()
	mockUser := models.UserDTO{ID: primitive.NewObjectID(), Username: "john"}
	mockUC.On("GetUserById", mock.Anything, userId).Return(mockUser, nil)

	r := httptest.NewRequest(http.MethodGet, "/user?userId="+userId, nil)
	w := httptest.NewRecorder()

	handler.GetUser(w, r)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	mockUC := new(mockUserUsecase)
	mockImageUC := new(mockImageUsecase)
	handler := NewUserHandler(mockUC, mockImageUC)

	userId := primitive.NewObjectID()
	userDTO := models.UserDTO{
		ID:       userId,
		Username: "john",
		Bio:      new(string),
		Email:    new(string),
	}
	*userDTO.Bio = "new bio"
	*userDTO.Email = "new@example.com"
	updatedToken := "newAccessToken"

	// Mock GetUserById to return the userDTO
	mockUC.On("GetUserById", mock.Anything, userId.Hex()).Return(userDTO, nil)

	// Mock UpdateUser to return the updated token
	mockUC.On("UpdateUser", mock.Anything, userDTO).Return(updatedToken, nil)

	// Prepare the request with multipart data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("bio", *userDTO.Bio)
	writer.WriteField("email", *userDTO.Email)
	writer.Close()

	r := httptest.NewRequest(http.MethodPut, "/user?userId="+userId.Hex(), body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// Call the handler
	handler.UpdateUser(w, r)

	// Check the response status
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	var resp map[string]string
	_ = json.NewDecoder(res.Body).Decode(&resp)
	if resp["access_token"] != updatedToken {
		t.Errorf("expected token %s, got %s", updatedToken, resp["access_token"])
	}
}
