package repository

import (
	"database/sql"
	"time"

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
	query := `INSERT INTO posts (author_id, category_id, title, content) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, post.AuthorID, post.CategoryID, post.Title, post.Content)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	post.ID = id
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	return nil
}

func (r *PostRepository) GetAll() ([]models.Post, error) {
	rows, err := r.db.Query(`SELECT id, author_id, category_id, title, content, created_at, updated_at FROM posts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.CategoryID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetByID(id int) (*models.Post, error) {
	var p models.Post
	err := r.db.QueryRow(`SELECT id, author_id, category_id, title, content, created_at, updated_at FROM posts WHERE id = ?`, id).Scan(
		&p.ID, &p.AuthorID, &p.CategoryID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostRepository) Update(post *models.Post) error {
	query := `UPDATE posts SET title = ?, content = ?, category_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, post.Title, post.Content, post.CategoryID, post.ID)
	return err
}

func (r *PostRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM posts WHERE id = ?`, id)
	return err
}
