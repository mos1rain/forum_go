package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/mos1rain/forum_go/internal/auth/models"
	"github.com/mos1rain/forum_go/internal/auth/service"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	service service.UserServiceInterface
	logger  zerolog.Logger
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger(),
	}
}

// @Summary Register new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.CreateUserInput true "User registration data"
// @Success 201 {object} service.AuthResponse "User successfully registered"
// @Failure 400 {string} string "Invalid request data"
// @Failure 409 {string} string "User already exists"
// @Failure 500 {string} string "Internal server error"
// @Router /api/auth/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Неверный формат данных"})
		return
	}

	h.logger.Info().Str("username", input.Username).Str("email", input.Email).Msg("Attempting to register new user")

	response, err := h.service.Register(input)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to register user")
		w.Header().Set("Content-Type", "application/json")
		switch err {
		case service.ErrUserAlreadyExists:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь с такими данными уже существует"})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Внутренняя ошибка сервера"})
		}
		return
	}

	h.logger.Info().Int("user_id", response.User.ID).Msg("User successfully registered")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
	}
}

// @Summary User login
// @Description Authenticate user and get access token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.LoginInput true "Login credentials"
// @Success 200 {object} service.AuthResponse "Login successful"
// @Failure 400 {string} string "Invalid request data"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Internal server error"
// @Router /api/auth/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Info().Str("username", input.Username).Msg("Attempting to login user")

	response, err := h.service.Login(input)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to login user")
		switch err {
		case service.ErrInvalidCredentials:
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info().Int("user_id", response.User.ID).Msg("User successfully logged in")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
	}
}
