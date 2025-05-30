package service

import (
	"errors"
	"time"

	"github.com/mos1rain/forum_go/internal/auth/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidInput       = errors.New("invalid input")
)

type UserRepo interface {
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByID(id int) (*models.User, error)
	Create(user *models.User) error
}

type TokenManager interface {
	NewJWTWithRole(userID int, username, role string, ttl time.Duration) (string, error)
}

type UserServiceInterface interface {
	Register(input models.CreateUserInput) (*AuthResponse, error)
	Login(input models.LoginInput) (*AuthResponse, error)
}

type UserService struct {
	repo         UserRepo
	tokenManager TokenManager
	tokenTTL     time.Duration
}

func NewUserService(repo UserRepo, tokenManager TokenManager, tokenTTL time.Duration) *UserService {
	return &UserService{
		repo:         repo,
		tokenManager: tokenManager,
		tokenTTL:     tokenTTL,
	}
}

type AuthResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func (s *UserService) Register(input models.CreateUserInput) (*AuthResponse, error) {
	// Проверяем, существует ли пользователь с таким username
	if user, _ := s.repo.GetByUsername(input.Username); user != nil {
		return nil, ErrUserAlreadyExists
	}

	// Проверяем, существует ли пользователь с таким email
	if user, _ := s.repo.GetByEmail(input.Email); user != nil {
		return nil, ErrUserAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := input.Role
	if role == "" {
		role = "user"
	}

	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Генерируем токен с ролью
	token, err := s.tokenManager.NewJWTWithRole(user.ID, user.Username, user.Role, s.tokenTTL)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *UserService) Login(input models.LoginInput) (*AuthResponse, error) {
	user, err := s.repo.GetByUsername(input.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Генерируем токен с ролью
	token, err := s.tokenManager.NewJWTWithRole(user.ID, user.Username, user.Role, s.tokenTTL)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}
