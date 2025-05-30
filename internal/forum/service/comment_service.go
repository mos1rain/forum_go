package service

import (
	"github.com/mos1rain/forum_go/internal/forum/models"
	"github.com/mos1rain/forum_go/internal/forum/repository"
)

type CommentService struct {
	repo repository.CommentRepositoryInterface
}

func (s *CommentService) Create(comment *models.Comment) error {
	return s.repo.Create(comment)
}
func (s *CommentService) GetByPostID(postID int) ([]models.Comment, error) {
	return s.repo.GetByPostID(postID)
}
func (s *CommentService) Delete(id int) error {
	return s.repo.Delete(id)
}
