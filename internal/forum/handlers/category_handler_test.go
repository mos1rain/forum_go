package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mos1rain/forum_go/internal/forum/models"
	"github.com/mos1rain/forum_go/internal/forum/service"
)

// MockCategoryService - мок для CategoryServiceInterface
type MockCategoryService struct {
	createFunc  func(ctx context.Context, category *models.Category) error
	getAllFunc  func(ctx context.Context) ([]*models.Category, error)
	getByIDFunc func(ctx context.Context, id int64) (*models.Category, error)
	deleteFunc  func(ctx context.Context, id int64, role string) error
}

func (m *MockCategoryService) Create(ctx context.Context, category *models.Category) error {
	return m.createFunc(ctx, category)
}

func (m *MockCategoryService) GetAll(ctx context.Context) ([]*models.Category, error) {
	return m.getAllFunc(ctx)
}

func (m *MockCategoryService) GetByID(ctx context.Context, id int64) (*models.Category, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *MockCategoryService) Delete(ctx context.Context, id int64, role string) error {
	return m.deleteFunc(ctx, id, role)
}

func TestCategoryHandler_DeleteCategory(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		categoryID     string
		userRole       string
		mockDeleteFunc func(ctx context.Context, id int64, role string) error
		wantStatus     int
	}{
		{
			name:       "successful deletion",
			method:     http.MethodDelete,
			categoryID: "1",
			userRole:   "admin",
			mockDeleteFunc: func(ctx context.Context, id int64, role string) error {
				return nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not admin role",
			method:     http.MethodDelete,
			categoryID: "1",
			userRole:   "user",
			mockDeleteFunc: func(ctx context.Context, id int64, role string) error {
				return service.ErrAdminRoleRequired
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "category not found",
			method:     http.MethodDelete,
			categoryID: "999",
			userRole:   "admin",
			mockDeleteFunc: func(ctx context.Context, id int64, role string) error {
				return service.ErrCategoryNotFound
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid category ID",
			method:     http.MethodDelete,
			categoryID: "invalid",
			userRole:   "admin",
			mockDeleteFunc: func(ctx context.Context, id int64, role string) error {
				return nil
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing category ID",
			method:     http.MethodDelete,
			categoryID: "",
			userRole:   "admin",
			mockDeleteFunc: func(ctx context.Context, id int64, role string) error {
				return nil
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "wrong HTTP method",
			method:     http.MethodGet,
			categoryID: "1",
			userRole:   "admin",
			mockDeleteFunc: func(ctx context.Context, id int64, role string) error {
				return nil
			},
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "unauthorized",
			method:     http.MethodDelete,
			categoryID: "1",
			userRole:   "",
			mockDeleteFunc: func(ctx context.Context, id int64, role string) error {
				return nil
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок сервиса
			mockService := &MockCategoryService{
				deleteFunc: tt.mockDeleteFunc,
			}
			handler := NewCategoryHandler(mockService)

			// Создаем тестовый запрос
			req := httptest.NewRequest(tt.method, "/api/categories?id="+tt.categoryID, nil)

			// Если роль указана, добавляем её в контекст
			if tt.userRole != "" {
				ctx := context.WithValue(req.Context(), "user_role", tt.userRole)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()

			// Выполняем запрос
			handler.DeleteCategory(rec, req)

			// Проверяем статус ответа
			if rec.Code != tt.wantStatus {
				t.Errorf("DeleteCategory() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}
