package auth

import (
	"ChatApp/internal/auth/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUsecase struct {
	signUpFunc func(ctx context.Context, user models.User) error
	signInFunc func(ctx context.Context, username, password string) (string, error)
}

func (m *mockUsecase) SignUp(ctx context.Context, user models.User) error {
	return m.signUpFunc(ctx, user)
}

func (m *mockUsecase) SignIn(ctx context.Context, username, password string) (string, error) {
	return m.signInFunc(ctx, username, password)
}

func TestSignUpHandler_Success(t *testing.T) {
	mock := &mockUsecase{
		signUpFunc: func(ctx context.Context, user models.User) error {
			return nil
		},
	}

	handler := NewAuthHandler(mock)

	user := models.UserCreateDTO{
		Username: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.SignUp(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}
}

func TestSignInHandler_Success(t *testing.T) {
	mock := &mockUsecase{
		signInFunc: func(ctx context.Context, username, password string) (string, error) {
			return "mocked-token", nil
		},
	}

	handler := NewAuthHandler(mock)

	user := models.UserCreateDTO{
		Username: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.SignIn(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var data map[string]string
	json.NewDecoder(res.Body).Decode(&data)
	if data["access_token"] != "mocked-token" {
		t.Errorf("Expected token 'mocked-token', got %s", data["access_token"])
	}
}
