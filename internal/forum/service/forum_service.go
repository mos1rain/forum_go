package service

import (
	"github.com/mos1rain/forum_go/internal/forum/models"
	"github.com/mos1rain/forum_go/internal/forum/repository"
)

type ForumService struct {
	Categories *CategoryService
	Posts      *PostService
	Comments   *CommentService
}

type CategoryService struct {
	repo repository.CategoryRepositoryInterface
}

type PostService struct {
	repo repository.PostRepositoryInterface
}

type CommentService struct {
	repo repository.CommentRepositoryInterface
}

func NewForumService(catRepo repository.CategoryRepositoryInterface, postRepo repository.PostRepositoryInterface, commRepo repository.CommentRepositoryInterface) *ForumService {
	return &ForumService{
		Categories: &CategoryService{repo: catRepo},
		Posts:      &PostService{repo: postRepo},
		Comments:   &CommentService{repo: commRepo},
	}
}

// Category methods
func (s *CategoryService) Create(category *models.Category) error {
	return s.repo.Create(category)
}
func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}
func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}
func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}

// Post methods
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

// Comment methods
func (s *CommentService) Create(comment *models.Comment) error {
	return s.repo.Create(comment)
}
func (s *CommentService) GetByPostID(postID int) ([]models.Comment, error) {
	return s.repo.GetByPostID(postID)
}
func (s *CommentService) Delete(id int) error {
	return s.repo.Delete(id)
}
