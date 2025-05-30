package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mos1rain/forum_go/docs"
	"github.com/mos1rain/forum_go/internal/auth/grpc"
	"github.com/mos1rain/forum_go/internal/auth/handler"
	"github.com/mos1rain/forum_go/internal/auth/repository"
	"github.com/mos1rain/forum_go/internal/auth/service"
	"github.com/mos1rain/forum_go/pkg/jwt"
	"github.com/rs/zerolog"
	_ "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type Category struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatorID   int64  `json:"creatorId"`
}

type Post struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CategoryID int64     `json:"categoryId"`
	AuthorID   int64     `json:"authorId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

var categories []Category
var posts []Post

// Обёртка для CORS
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func checkAdminExists(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func createDefaultAdmin(db *sql.DB, logger zerolog.Logger) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "admin", "admin@forum.com", string(hashedPassword), "admin", time.Now(), time.Now())

	if err != nil {
		return err
	}

	logger.Info().
		Str("username", "admin").
		Str("email", "admin@forum.com").
		Msg("Default admin user created")

	return nil
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	// Подключение к SQLite с правильными настройками
	db, err := sql.Open("sqlite", "/Users/Sieger/Desktop/forum_go/forum.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to ping database")
	}

	logger.Info().Msg("Successfully connected to database")

	// Инициализация таблиц
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize users table")
	}

	// Проверяем наличие администратора и создаем его, если нет
	adminExists, err := checkAdminExists(db)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to check admin existence")
	}

	if !adminExists {
		if err := createDefaultAdmin(db, logger); err != nil {
			logger.Fatal().Err(err).Msg("Failed to create default admin")
		}
	}

	logger.Info().Msg("Database tables initialized successfully")

	// Инициализируем JWT менеджер
	tokenManager := jwt.NewTokenManager(jwt.SecretKey)
	tokenTTL := 24 * time.Hour

	// Инициализируем слои приложения
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, tokenManager, tokenTTL)
	userHandler := handler.NewUserHandler(userService)

	// Запуск gRPC-сервера в отдельной горутине
	go grpc.RunGRPCServer(userRepo, tokenManager, ":50052")

	// Регистрируем маршруты с CORS
	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/register", withCORS(userHandler.Register))
	mux.HandleFunc("/api/auth/login", withCORS(userHandler.Login))
	mux.HandleFunc("/api/categories", withCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(categories)
		case http.MethodPost:
			// Получаем ID пользователя из токена
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			token = strings.TrimPrefix(token, "Bearer ")

			// Декодируем токен и получаем ID пользователя
			claims, err := tokenManager.Parse(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			logger.Info().
				Int("user_id", claims.UserID).
				Str("username", claims.Username).
				Msg("Creating new category")

			var category Category
			if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Устанавливаем ID создателя категории
			category.ID = int64(len(categories) + 1)
			category.CreatorID = int64(claims.UserID)

			logger.Info().
				Int64("category_id", category.ID).
				Int64("creator_id", category.CreatorID).
				Str("name", category.Name).
				Msg("New category created")

			categories = append(categories, category)
			json.NewEncoder(w).Encode(category)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	mux.HandleFunc("/api/categories/", withCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Извлекаем ID категории из URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		categoryID, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		// Проверяем, существует ли категория
		var categoryExists bool
		for _, cat := range categories {
			if cat.ID == categoryID {
				categoryExists = true
				break
			}
		}

		if !categoryExists {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		// Если запрос к /api/categories/{id}/posts
		if len(parts) > 4 && parts[4] == "posts" {
			// Фильтруем посты по категории
			var categoryPosts []Post
			for _, post := range posts {
				if post.CategoryID == categoryID {
					categoryPosts = append(categoryPosts, post)
				}
			}
			json.NewEncoder(w).Encode(categoryPosts)
			return
		}

		http.Error(w, "Not found", http.StatusNotFound)
	}))
	mux.HandleFunc("/api/posts", withCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(posts)
		case http.MethodPost:
			var post Post
			if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			post.ID = int64(len(posts) + 1)
			post.CreatedAt = time.Now()
			post.UpdatedAt = time.Now()
			posts = append(posts, post)
			json.NewEncoder(w).Encode(post)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	mux.HandleFunc("/api/forum/delete_category", withCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получаем ID пользователя из токена
		token := r.Header.Get("Authorization")
		if token == "" {
			logger.Error().Msg("No authorization token provided")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		logger.Info().
			Str("token", token).
			Msg("Received token for category deletion")

		// Декодируем токен и получаем ID пользователя
		claims, err := tokenManager.Parse(token)
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to parse token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		logger.Info().
			Int("current_user_id", claims.UserID).
			Str("username", claims.Username).
			Str("token", token).
			Msg("Successfully parsed token")

		// Выводим информацию о всех существующих категориях
		logger.Info().Msg("Current categories:")
		for _, cat := range categories {
			logger.Info().
				Int64("category_id", cat.ID).
				Int64("creator_id", cat.CreatorID).
				Str("name", cat.Name).
				Msg("Category info")
		}

		idStr := r.URL.Query().Get("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			logger.Error().
				Err(err).
				Str("id", idStr).
				Msg("Invalid category ID format")
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		logger.Info().
			Int64("requested_category_id", id).
			Int("current_user_id", claims.UserID).
			Msg("Attempting to delete category")

		// Находим категорию и проверяем права доступа
		for i, cat := range categories {
			if cat.ID == id {
				logger.Info().
					Int64("category_id", cat.ID).
					Int64("category_creator_id", cat.CreatorID).
					Int("current_user_id", claims.UserID).
					Str("category_name", cat.Name).
					Msg("Found category, checking permissions")

				// Проверяем, является ли пользователь создателем категории
				if cat.CreatorID != int64(claims.UserID) {
					logger.Warn().
						Int64("category_id", cat.ID).
						Int64("category_creator_id", cat.CreatorID).
						Int("current_user_id", claims.UserID).
						Str("category_name", cat.Name).
						Msg("Permission denied: user is not the category creator")
					http.Error(w, "You don't have permission to delete this category", http.StatusForbidden)
					return
				}

				categories = append(categories[:i], categories[i+1:]...)
				logger.Info().
					Int64("category_id", cat.ID).
					Int64("creator_id", cat.CreatorID).
					Str("category_name", cat.Name).
					Msg("Category successfully deleted")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully"})
				return
			}
		}

		logger.Warn().
			Int64("requested_category_id", id).
			Msg("Category not found")
		http.Error(w, "Category not found", http.StatusNotFound)
	}))
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Запускаем HTTP сервер
	logger.Info().Str("port", "3001").Msg("Starting auth server on :3001")
	if err := http.ListenAndServe(":3001", mux); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
