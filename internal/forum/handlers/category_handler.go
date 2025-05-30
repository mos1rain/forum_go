package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/mos1rain/forum_go/internal/forum/models"
	"github.com/mos1rain/forum_go/internal/forum/service"
)

var (
	ErrAdminRoleRequired = errors.New("only admin can delete categories")
)

type CategoryServiceInterface interface {
	Create(ctx context.Context, category *models.Category) error
	GetAll(ctx context.Context) ([]*models.Category, error)
	GetByID(ctx context.Context, id int64) (*models.Category, error)
	Delete(ctx context.Context, id int64, role string) error
}

type CategoryHandler struct {
	service CategoryServiceInterface
}

func NewCategoryHandler(service CategoryServiceInterface) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID категории из query параметров
	categoryID := r.URL.Query().Get("id")
	if categoryID == "" {
		http.Error(w, "category ID is required", http.StatusBadRequest)
		return
	}

	// Получаем роль пользователя из контекста
	userRole, ok := r.Context().Value("user_role").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(categoryID, 10, 64)
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	// Удаляем категорию с проверкой прав
	err = h.service.Delete(r.Context(), id, userRole)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAdminRoleRequired):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, service.ErrCategoryNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
