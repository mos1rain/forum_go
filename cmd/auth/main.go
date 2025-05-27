package main

import (
	"net/http"
	"os"
	"time"

	_ "github.com/mos1rain/forum_go/docs"
	"github.com/mos1rain/forum_go/internal/auth/grpc"
	"github.com/mos1rain/forum_go/internal/auth/handler"
	"github.com/mos1rain/forum_go/internal/auth/repository"
	"github.com/mos1rain/forum_go/internal/auth/service"
	"github.com/mos1rain/forum_go/pkg/database"
	"github.com/mos1rain/forum_go/pkg/jwt"
	"github.com/rs/zerolog"
	_ "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	// Подключаемся к базе данных
	db, err := database.NewPostgresDB(
		"localhost",
		"5432",
		"postgres",
		"28072005",
		"forum",
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Инициализируем JWT менеджер
	tokenManager := jwt.NewTokenManager("your-secret-key") // В продакшене использовать безопасный ключ
	tokenTTL := 24 * time.Hour

	// Инициализируем слои приложения
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, tokenManager, tokenTTL)
	userHandler := handler.NewUserHandler(userService)

	// Запуск gRPC-сервера в отдельной горутине
	go grpc.RunGRPCServer(userRepo, tokenManager, ":50051")

	// Регистрируем маршруты
	http.HandleFunc("/api/auth/register", withCORS(userHandler.Register))
	http.HandleFunc("/api/auth/login", withCORS(userHandler.Login))
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Запускаем сервер
	logger.Info().Msg("Starting auth server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}
