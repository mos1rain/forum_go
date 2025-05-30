package service

import (
	"github.com/mos1rain/forum_go/internal/forum/models"
	"github.com/mos1rain/forum_go/internal/forum/repository"
)

type PostService struct {
	repo repository.PostRepositoryInterface
}

func (s *PostService) Create(post *models.Post) error {
	return s.repo.Create(post)
}
func (s *PostService) GetAll() ([]models.Post, error) {
	return s.repo.GetAll()
}
func (s *PostService) GetByID(id int) (*models.Post, error) {
	return s.repo.GetByID(id)
}
func (s *PostService) Update(post *models.Post) error {
	return s.repo.Update(post)
}
func (s *PostService) Delete(id int) error {
	return s.repo.Delete(id)
}
