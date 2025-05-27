package service

import (
	"errors"
	"testing"
	"time"

	"github.com/mos1rain/forum_go/internal/auth/models"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	users map[string]*models.User
}

var _ UserRepo = (*mockUserRepo)(nil)

func (m *mockUserRepo) GetByUsername(username string) (*models.User, error) {
	if u, ok := m.users[username]; ok {
		return u, nil
	}
	return nil, nil
}
func (m *mockUserRepo) GetByEmail(email string) (*models.User, error) { return nil, nil }
func (m *mockUserRepo) GetByID(id int) (*models.User, error)          { return nil, nil }
func (m *mockUserRepo) Create(user *models.User) error {
	if _, ok := m.users[user.Username]; ok {
		return errors.New("already exists")
	}
	m.users[user.Username] = user
	return nil
}

func TestRegister_NewUser(t *testing.T) {
	repo := &mockUserRepo{users: map[string]*models.User{}}
	tm := newTestTokenManager()
	s := NewUserService(repo, tm, 0)
	input := models.CreateUserInput{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	resp, err := s.Register(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.User.Username != input.Username {
		t.Errorf("expected username %s, got %s", input.Username, resp.User.Username)
	}
}

func TestRegister_AlreadyExists(t *testing.T) {
	repo := &mockUserRepo{users: map[string]*models.User{"testuser": {Username: "testuser"}}}
	tm := newTestTokenManager()
	s := NewUserService(repo, tm, 0)
	input := models.CreateUserInput{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	_, err := s.Register(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLogin_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	repo := &mockUserRepo{users: map[string]*models.User{
		"testuser": {ID: 1, Username: "testuser", PasswordHash: string(hash), Role: "user"},
	}}
	tm := newTestTokenManager()
	s := NewUserService(repo, tm, 0)
	input := models.LoginInput{
		Username: "testuser",
		Password: "password",
	}
	resp, err := s.Login(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.User.Username != input.Username {
		t.Errorf("expected username %s, got %s", input.Username, resp.User.Username)
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	repo := &mockUserRepo{users: map[string]*models.User{
		"testuser": {ID: 1, Username: "testuser", PasswordHash: string(hash), Role: "user"},
	}}
	tm := newTestTokenManager()
	s := NewUserService(repo, tm, 0)
	input := models.LoginInput{
		Username: "testuser",
		Password: "wrongpass",
	}
	_, err := s.Login(input)
	if err == nil {
		t.Fatal("expected error for invalid password")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{users: map[string]*models.User{}}
	tm := newTestTokenManager()
	s := NewUserService(repo, tm, 0)
	input := models.LoginInput{
		Username: "nouser",
		Password: "password",
	}
	_, err := s.Login(input)
	if err == nil {
		t.Fatal("expected error for user not found")
	}
}

type errorRepo struct {
	mockUserRepo
}

func (e *errorRepo) GetByUsername(username string) (*models.User, error) {
	return nil, errors.New("repo error")
}

func TestLogin_RepoError(t *testing.T) {
	repo := &errorRepo{mockUserRepo{users: map[string]*models.User{}}}
	tm := newTestTokenManager()
	s := NewUserService(repo, tm, 0)
	input := models.LoginInput{
		Username: "testuser",
		Password: "password",
	}
	_, err := s.Login(input)
	if err == nil || err.Error() != "repo error" {
		t.Fatalf("expected repo error, got %v", err)
	}
}

type testTokenManager struct{}

var _ TokenManager = (*testTokenManager)(nil)

func newTestTokenManager() *testTokenManager { return &testTokenManager{} }
func (t *testTokenManager) NewJWTWithRole(userID int, username, role string, ttl time.Duration) (string, error) {
	return "token", nil
}
func (t *testTokenManager) NewJWT(userID int, username string, ttl int64) (string, error) {
	return "token", nil
}
