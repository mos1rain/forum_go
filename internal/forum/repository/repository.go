package repository

import (
	"context"

	"github.com/mos1rain/forum_go/internal/forum/models"
)

// Repository предоставляет методы для работы с базой данных
type Repository interface {
	// CreateCategory создает новую категорию в базе данных
	CreateCategory(ctx context.Context, category *models.Category) error
	// DeleteCategory удаляет категорию из базы данных
	DeleteCategory(ctx context.Context, id int64) error
	// GetCategories возвращает список всех категорий из базы данных
	GetCategories(ctx context.Context) ([]*models.Category, error)
	// GetCategoryByID возвращает категорию по ID
	GetCategoryByID(ctx context.Context, id int64) (*models.Category, error)

	// CreatePost создает новый пост в базе данных
	CreatePost(ctx context.Context, post *models.Post) error
	// DeletePost удаляет пост из базы данных
	DeletePost(ctx context.Context, id int64) error
	// GetPosts возвращает список постов в категории
	GetPosts(ctx context.Context, categoryID int64) ([]*models.Post, error)
	// GetPostByID возвращает пост по ID
	GetPostByID(ctx context.Context, id int64) (*models.Post, error)

	// CreateComment создает новый комментарий в базе данных
	CreateComment(ctx context.Context, comment *models.Comment) error
	// DeleteComment удаляет комментарий из базы данных
	DeleteComment(ctx context.Context, id int64) error
	// GetComments возвращает список комментариев к посту
	GetComments(ctx context.Context, postID int64) ([]*models.Comment, error)
}
