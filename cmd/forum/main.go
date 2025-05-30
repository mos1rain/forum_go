package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mos1rain/forum_go/docs"
	"github.com/mos1rain/forum_go/internal/forum/grpc"
	"github.com/mos1rain/forum_go/internal/forum/handler"
	"github.com/mos1rain/forum_go/internal/forum/middleware"
	"github.com/mos1rain/forum_go/internal/forum/repository"
	"github.com/mos1rain/forum_go/internal/forum/service"
	"github.com/mos1rain/forum_go/pkg/jwt"
	"github.com/rs/zerolog"
	_ "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "modernc.org/sqlite"
)

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
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			description TEXT NOT NULL,
			creator_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize database tables")
	}

	logger.Info().Msg("Database tables initialized successfully")

	// Инициализация gRPC клиента для аутентификации
	authClient, err := grpc.NewAuthGRPCClient("localhost:50052")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to auth service")
	}
	middleware.SetAuthClient(authClient)

	catRepo := repository.NewCategoryRepository(db)
	postRepo := repository.NewPostRepository(db)
	commRepo := repository.NewCommentRepository(db)
	forumService := service.NewForumService(catRepo, postRepo, commRepo)
	h := handler.NewForumHandler(forumService)

	// Создаем TokenManager с тем же секретным ключом
	tokenManager := jwt.NewTokenManager(jwt.SecretKey)
	middleware.SetTokenManager(tokenManager)

	// Создаем новый маршрутизатор
	mux := http.NewServeMux()

	// Регистрируем маршруты с CORS
	mux.HandleFunc("/api/forum/categories", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetCategories(w, r)
		} else if r.Method == http.MethodPost {
			middleware.AuthMiddleware(http.HandlerFunc(h.CreateCategory)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/forum/posts", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetPosts(w, r)
		} else if r.Method == http.MethodPost {
			middleware.AuthMiddleware(http.HandlerFunc(h.CreatePost)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/forum/comments", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetCommentsByPost(w, r)
		} else if r.Method == http.MethodPost {
			middleware.AuthMiddleware(http.HandlerFunc(h.CreateComment)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/forum/posts/", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			idStr := r.URL.Path[len("/api/forum/posts/"):]
			if idStr == "" {
				http.Error(w, "Missing post id", http.StatusBadRequest)
				return
			}
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid post id", http.StatusBadRequest)
				return
			}
			post, err := h.GetPostByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(post)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}))

	mux.HandleFunc("/api/forum/delete_post", withCORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.AuthMiddleware(http.HandlerFunc(h.DeletePost)).ServeHTTP(w, r)
	}))

	mux.HandleFunc("/api/forum/delete_comment", withCORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.AuthMiddleware(http.HandlerFunc(h.DeleteComment)).ServeHTTP(w, r)
	}))

	mux.HandleFunc("/api/forum/delete_category", withCORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.AuthMiddleware(http.HandlerFunc(h.DeleteCategory)).ServeHTTP(w, r)
	}))

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	logger.Info().Str("port", "3002").Msg("Starting forum server on :3002")
	if err := http.ListenAndServe(":3002", mux); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
