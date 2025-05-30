package repository

import (
	"database/sql"
	"time"

	"github.com/mos1rain/forum_go/internal/forum/models"
)

type CommentRepository struct {
	db *sql.DB
}

type CommentRepositoryInterface interface {
	Create(comment *models.Comment) error
	GetByPostID(postID int) ([]models.Comment, error)
	Delete(id int) error
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, comment.PostID, comment.AuthorID, comment.Content)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	comment.ID = id
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	return nil
}

func (r *CommentRepository) GetByPostID(postID int) ([]models.Comment, error) {
	rows, err := r.db.Query(`SELECT id, post_id, user_id, content, created_at, updated_at FROM comments WHERE post_id = ?`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.AuthorID, &c.Content, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *CommentRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM comments WHERE id = ?`, id)
	return err
}
