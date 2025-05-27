package service

import (
	"errors"
	"testing"

	"github.com/mos1rain/forum_go/internal/forum/models"
	"github.com/mos1rain/forum_go/internal/forum/repository"
)

type mockCategoryRepo struct{ cats []models.Category }

var _ repository.CategoryRepositoryInterface = (*mockCategoryRepo)(nil)

func (m *mockCategoryRepo) Create(cat *models.Category) error {
	m.cats = append(m.cats, *cat)
	return nil
}
func (m *mockCategoryRepo) GetAll() ([]models.Category, error) { return m.cats, nil }
func (m *mockCategoryRepo) Delete(id int) error                { return nil }
func (m *mockCategoryRepo) GetByID(id int) (*models.Category, error) {
	for _, c := range m.cats {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, errors.New("not found")
}

type mockPostRepo struct{ posts []models.Post }

var _ repository.PostRepositoryInterface = (*mockPostRepo)(nil)

func (m *mockPostRepo) Create(post *models.Post) error {
	m.posts = append(m.posts, *post)
	return nil
}
func (m *mockPostRepo) GetAll() ([]models.Post, error) { return m.posts, nil }
func (m *mockPostRepo) GetByID(id int) (*models.Post, error) {
	for _, p := range m.posts {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, errors.New("not found")
}
func (m *mockPostRepo) Delete(id int) error { return nil }
func (m *mockPostRepo) Update(post *models.Post) error {
	for i, p := range m.posts {
		if p.ID == post.ID {
			m.posts[i] = *post
			return nil
		}
	}
	return errors.New("not found")
}

type mockCommentRepo struct{ comms []models.Comment }

var _ repository.CommentRepositoryInterface = (*mockCommentRepo)(nil)

func (m *mockCommentRepo) Create(c *models.Comment) error {
	m.comms = append(m.comms, *c)
	return nil
}
func (m *mockCommentRepo) GetByPostID(postID int) ([]models.Comment, error) {
	var res []models.Comment
	for _, c := range m.comms {
		if c.PostID == postID {
			res = append(res, c)
		}
	}
	return res, nil
}
func (m *mockCommentRepo) Delete(id int) error { return nil }

func TestCreateAndGetCategory(t *testing.T) {
	catRepo := &mockCategoryRepo{}
	fs := NewForumService(catRepo, &mockPostRepo{}, &mockCommentRepo{})
	cat := &models.Category{Name: "TestCat", Description: "desc"}
	if err := fs.Categories.Create(cat); err != nil {
		t.Fatalf("create: %v", err)
	}
	cats, err := fs.Categories.GetAll()
	if err != nil || len(cats) != 1 {
		t.Fatalf("get all: %v", err)
	}
}

func TestCreateAndGetPost(t *testing.T) {
	postRepo := &mockPostRepo{}
	fs := NewForumService(&mockCategoryRepo{}, postRepo, &mockCommentRepo{})
	post := &models.Post{ID: 1, Title: "Test", Content: "Body", CategoryID: 1, UserID: 1}
	if err := fs.Posts.Create(post); err != nil {
		t.Fatalf("create: %v", err)
	}
	posts, err := fs.Posts.GetAll()
	if err != nil || len(posts) != 1 {
		t.Fatalf("get all: %v", err)
	}
}

func TestCreateAndGetComment(t *testing.T) {
	commRepo := &mockCommentRepo{}
	fs := NewForumService(&mockCategoryRepo{}, &mockPostRepo{}, commRepo)
	comm := &models.Comment{ID: 1, PostID: 1, Content: "Test comment", UserID: 1}
	if err := fs.Comments.Create(comm); err != nil {
		t.Fatalf("create: %v", err)
	}
	comms, err := fs.Comments.GetByPostID(1)
	if err != nil || len(comms) != 1 {
		t.Fatalf("get by post: %v", err)
	}
}

func TestCategoryGetByID(t *testing.T) {
	catRepo := &mockCategoryRepo{
		cats: []models.Category{
			{ID: 1, Name: "TestCat1", Description: "desc1"},
			{ID: 2, Name: "TestCat2", Description: "desc2"},
		},
	}
	fs := NewForumService(catRepo, &mockPostRepo{}, &mockCommentRepo{})

	// Test existing category
	cat, err := fs.Categories.GetByID(1)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if cat.Name != "TestCat1" {
		t.Errorf("expected TestCat1, got %s", cat.Name)
	}

	// Test non-existing category
	cat, err = fs.Categories.GetByID(999)
	if err == nil {
		t.Error("expected error for non-existing category")
	}
}

func TestCategoryDelete(t *testing.T) {
	catRepo := &mockCategoryRepo{
		cats: []models.Category{
			{ID: 1, Name: "TestCat", Description: "desc"},
		},
	}
	fs := NewForumService(catRepo, &mockPostRepo{}, &mockCommentRepo{})

	if err := fs.Categories.Delete(1); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestPostUpdate(t *testing.T) {
	postRepo := &mockPostRepo{
		posts: []models.Post{
			{ID: 1, Title: "Old Title", Content: "Old Content", CategoryID: 1, UserID: 1},
		},
	}
	fs := NewForumService(&mockCategoryRepo{}, postRepo, &mockCommentRepo{})

	updatedPost := &models.Post{
		ID:         1,
		Title:      "New Title",
		Content:    "New Content",
		CategoryID: 2,
		UserID:     1,
	}

	if err := fs.Posts.Update(updatedPost); err != nil {
		t.Fatalf("update: %v", err)
	}

	// Verify update
	post, err := fs.Posts.GetByID(1)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if post.Title != "New Title" {
		t.Errorf("expected New Title, got %s", post.Title)
	}
	if post.Content != "New Content" {
		t.Errorf("expected New Content, got %s", post.Content)
	}
	if post.CategoryID != 2 {
		t.Errorf("expected category ID 2, got %d", post.CategoryID)
	}
}

func TestPostUpdateNonExisting(t *testing.T) {
	postRepo := &mockPostRepo{}
	fs := NewForumService(&mockCategoryRepo{}, postRepo, &mockCommentRepo{})

	nonExistingPost := &models.Post{
		ID:         999,
		Title:      "New Title",
		Content:    "New Content",
		CategoryID: 1,
		UserID:     1,
	}

	if err := fs.Posts.Update(nonExistingPost); err == nil {
		t.Error("expected error for updating non-existing post")
	}
}

func TestCommentDelete(t *testing.T) {
	commRepo := &mockCommentRepo{
		comms: []models.Comment{
			{ID: 1, PostID: 1, Content: "Test comment", UserID: 1},
		},
	}
	fs := NewForumService(&mockCategoryRepo{}, &mockPostRepo{}, commRepo)

	if err := fs.Comments.Delete(1); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
