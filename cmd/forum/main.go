package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	_ "github.com/mos1rain/forum_go/docs"
	"github.com/mos1rain/forum_go/internal/forum/handler"
	"github.com/mos1rain/forum_go/internal/forum/middleware"
	"github.com/mos1rain/forum_go/internal/forum/repository"
	"github.com/mos1rain/forum_go/internal/forum/service"
	"github.com/mos1rain/forum_go/pkg/jwt"
	"github.com/rs/zerolog"
	_ "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Обёртка для CORS
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
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

	db, err := sql.Open("postgres", "postgres://postgres:28072005@localhost:5432/forum?sslmode=disable")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	catRepo := repository.NewCategoryRepository(db)
	postRepo := repository.NewPostRepository(db)
	commRepo := repository.NewCommentRepository(db)
	forumService := service.NewForumService(catRepo, postRepo, commRepo)
	h := handler.NewForumHandler(forumService)

	tokenManager := jwt.NewTokenManager("your-secret-key") // тот же ключ, что и в Auth Service

	http.HandleFunc("/api/forum/categories", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetCategories(w, r)
		} else if r.Method == http.MethodPost {
			middleware.JWTAuth(tokenManager, http.HandlerFunc(h.CreateCategory)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/forum/posts", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetPosts(w, r)
		} else if r.Method == http.MethodPost {
			middleware.JWTAuth(tokenManager, http.HandlerFunc(h.CreatePost)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/forum/comments", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetCommentsByPost(w, r)
		} else if r.Method == http.MethodPost {
			middleware.JWTAuth(tokenManager, http.HandlerFunc(h.CreateComment)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/forum/posts/", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Получаем id из URL
			idStr := r.URL.Path[len("/api/forum/posts/"):]
			if idStr == "" {
				http.Error(w, "Missing post id", http.StatusBadRequest)
				return
			}
			// Преобразуем id к int
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

	http.HandleFunc("/api/forum/delete_post", withCORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuth(tokenManager, http.HandlerFunc(h.DeletePost)).ServeHTTP(w, r)
	}))

	http.HandleFunc("/api/forum/delete_comment", withCORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuth(tokenManager, http.HandlerFunc(h.DeleteComment)).ServeHTTP(w, r)
	}))

	http.HandleFunc("/api/forum/delete_category", withCORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuth(tokenManager, http.HandlerFunc(h.DeleteCategory)).ServeHTTP(w, r)
	}))

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	logger.Info().Msg("Forum service started on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start forum server")
	}
}
