package repository

import (
	"database/sql"

	"github.com/mos1rain/forum_go/internal/forum/models"
)

type PostRepository struct {
	db *sql.DB
}

type PostRepositoryInterface interface {
	Create(post *models.Post) error
	GetAll() ([]models.Post, error)
	GetByID(id int) (*models.Post, error)
	Update(post *models.Post) error
	Delete(id int) error
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post) error {
	query := `INSERT INTO posts (user_id, category_id, title, content) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query, post.UserID, post.CategoryID, post.Title, post.Content).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (r *PostRepository) GetAll() ([]models.Post, error) {
	rows, err := r.db.Query(`SELECT id, user_id, category_id, title, content, created_at, updated_at FROM posts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.CategoryID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetByID(id int) (*models.Post, error) {
	var p models.Post
	err := r.db.QueryRow(`SELECT id, user_id, category_id, title, content, created_at, updated_at FROM posts WHERE id = $1`, id).Scan(
		&p.ID, &p.UserID, &p.CategoryID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostRepository) Update(post *models.Post) error {
	query := `UPDATE posts SET title = $1, content = $2, category_id = $3, updated_at = NOW() WHERE id = $4`
	_, err := r.db.Exec(query, post.Title, post.Content, post.CategoryID, post.ID)
	return err
}

func (r *PostRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM posts WHERE id = $1`, id)
	return err
}
