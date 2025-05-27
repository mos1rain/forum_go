package repository

import (
	"database/sql"

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
	query := `INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(query, comment.PostID, comment.UserID, comment.Content).Scan(&comment.ID, &comment.CreatedAt)
}

func (r *CommentRepository) GetByPostID(postID int) ([]models.Comment, error) {
	rows, err := r.db.Query(`SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = $1`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *CommentRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM comments WHERE id = $1`, id)
	return err
}
