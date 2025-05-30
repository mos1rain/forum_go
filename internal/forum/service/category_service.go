package service

import (
	"context"
	"errors"

	"github.com/mos1rain/forum_go/internal/forum/models"
	"github.com/mos1rain/forum_go/internal/forum/repository"
)

var (
	ErrCategoryNotFound  = errors.New("category not found")
	ErrAdminRoleRequired = errors.New("only admin can delete categories")
)

type CategoryService struct {
	repo repository.CategoryRepositoryInterface
}

func NewCategoryService(repo repository.CategoryRepositoryInterface) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) Create(ctx context.Context, category *models.Category) error {
	return s.repo.CreateCategory(ctx, category)
}

func (s *CategoryService) GetAll(ctx context.Context) ([]*models.Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *CategoryService) GetByID(ctx context.Context, id int64) (*models.Category, error) {
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, id int64, role string) error {
	// Проверяем роль пользователя
	if role != "admin" {
		return ErrAdminRoleRequired
	}

	// Проверяем существование категории
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return ErrCategoryNotFound
	}

	// Удаляем категорию
	return s.repo.DeleteCategory(ctx, id)
}
