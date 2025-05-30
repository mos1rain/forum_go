package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/mos1rain/forum_go/internal/forum/models"
)

type CategoryRepository struct {
	db *sql.DB
}

type CategoryRepositoryInterface interface {
	CreateCategory(ctx context.Context, category *models.Category) error
	GetCategories(ctx context.Context) ([]*models.Category, error)
	GetCategoryByID(ctx context.Context, id int64) (*models.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	query := `INSERT INTO categories (name, description, creator_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, category.Name, category.Description, category.CreatorID, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	category.ID = id
	category.CreatedAt = now
	category.UpdatedAt = now

	return nil
}

func (r *CategoryRepository) GetCategories(ctx context.Context) ([]*models.Category, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, description, creator_id, created_at, updated_at FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatorID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id int64) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRowContext(ctx, `SELECT id, name, description, creator_id, created_at, updated_at FROM categories WHERE id = ?`, id).
		Scan(&c.ID, &c.Name, &c.Description, &c.CreatorID, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = ?`, id)
	return err
}
