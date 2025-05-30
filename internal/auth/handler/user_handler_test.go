package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mos1rain/forum_go/internal/auth/models"
	"github.com/mos1rain/forum_go/internal/auth/service"
)

type mockUserService struct {
	registerFunc func(input models.CreateUserInput) (*service.AuthResponse, error)
	loginFunc    func(input models.LoginInput) (*service.AuthResponse, error)
}

func (m *mockUserService) Register(input models.CreateUserInput) (*service.AuthResponse, error) {
	return m.registerFunc(input)
}

func (m *mockUserService) Login(input models.LoginInput) (*service.AuthResponse, error) {
	return m.loginFunc(input)
}

func TestUserHandler_Register(t *testing.T) {
	tests := []struct {
		name          string
		input         models.CreateUserInput
		mockRegister  func(input models.CreateUserInput) (*service.AuthResponse, error)
		expectedCode  int
		expectedError bool
	}{
		{
			name: "successful registration",
			input: models.CreateUserInput{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockRegister: func(input models.CreateUserInput) (*service.AuthResponse, error) {
				return &service.AuthResponse{
					User: &models.User{
						ID:       1,
						Username: input.Username,
						Email:    input.Email,
						Role:     "user",
					},
					Token: "test.jwt.token",
				}, nil
			},
			expectedCode:  http.StatusCreated,
			expectedError: false,
		},
		{
			name: "user already exists",
			input: models.CreateUserInput{
				Username: "existinguser",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockRegister: func(input models.CreateUserInput) (*service.AuthResponse, error) {
				return nil, service.ErrUserAlreadyExists
			},
			expectedCode:  http.StatusConflict,
			expectedError: true,
		},
		{
			name: "invalid input",
			input: models.CreateUserInput{
				Username: "",
				Email:    "invalid",
				Password: "",
			},
			mockRegister: func(input models.CreateUserInput) (*service.AuthResponse, error) {
				return nil, service.ErrInvalidInput
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем mock сервис
			mockService := &mockUserService{
				registerFunc: tt.mockRegister,
			}
			handler := NewUserHandler(mockService)

			// Создаем тестовый запрос
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			// Выполняем запрос
			handler.Register(rec, req)

			// Проверяем результат
			if rec.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, rec.Code)
			}

			if tt.expectedError {
				var errorResponse map[string]string
				if err := json.NewDecoder(rec.Body).Decode(&errorResponse); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if errorResponse["message"] == "" {
					t.Error("expected error message in response")
				}
			} else {
				var response service.AuthResponse
				if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode success response: %v", err)
				}
				if response.User == nil {
					t.Error("expected user in response")
				}
				if response.Token == "" {
					t.Error("expected token in response")
				}
			}
		})
	}
}

func TestUserHandler_Register_InvalidMethod(t *testing.T) {
	handler := NewUserHandler(&mockUserService{})
	req := httptest.NewRequest(http.MethodGet, "/api/auth/register", nil)
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status code %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestUserHandler_Register_InvalidJSON(t *testing.T) {
	handler := NewUserHandler(&mockUserService{})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
