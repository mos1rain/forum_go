package repository

import (
	"database/sql"

	"github.com/mos1rain/forum_go/internal/forum/models"
)

type CategoryRepository struct {
	db *sql.DB
}

type CategoryRepositoryInterface interface {
	Create(category *models.Category) error
	GetAll() ([]models.Category, error)
	GetByID(id int) (*models.Category, error)
	Delete(id int) error
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query(`SELECT id, name, description FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow(`SELECT id, name, description FROM categories WHERE id = $1`, id).Scan(&c.ID, &c.Name, &c.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	return err
}
