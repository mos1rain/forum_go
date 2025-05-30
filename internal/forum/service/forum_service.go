package service

import (
	"context"

	"github.com/mos1rain/forum_go/internal/forum/repository"
)

type ForumService struct {
	Categories *CategoryService
	Posts      *PostService
	Comments   *CommentService
}

func NewForumService(catRepo repository.CategoryRepositoryInterface, postRepo repository.PostRepositoryInterface, commRepo repository.CommentRepositoryInterface) *ForumService {
	return &ForumService{
		Categories: &CategoryService{repo: catRepo},
		Posts:      &PostService{repo: postRepo},
		Comments:   &CommentService{repo: commRepo},
	}
}

func (s *ForumService) DeleteCategory(id int, userRole string) error {
	return s.Categories.Delete(context.Background(), int64(id), userRole)
}
