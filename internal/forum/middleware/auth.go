package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/mos1rain/forum_go/internal/forum/grpc"
	"github.com/mos1rain/forum_go/pkg/jwt"
)

var (
	authClient   *grpc.AuthGRPCClient
	tokenManager *jwt.TokenManager
)

func SetAuthClient(client *grpc.AuthGRPCClient) {
	authClient = client
}

func SetTokenManager(tm *jwt.TokenManager) {
	tokenManager = tm
	fmt.Printf("TokenManager set with key: %v\n", tm.GetSigningKey())
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		fmt.Printf("Received token: %s\n", token)

		if token == "" {
			http.Error(w, "unauthorized: no token provided", http.StatusUnauthorized)
			return
		}

		// Убираем "Bearer " из токена
		token = strings.TrimPrefix(token, "Bearer ")
		fmt.Printf("Token after trim: %s\n", token)

		// Валидируем токен
		if tokenManager == nil {
			fmt.Println("TokenManager is nil!")
			http.Error(w, "token manager not initialized", http.StatusInternalServerError)
			return
		}

		claims, err := tokenManager.Parse(token)
		if err != nil {
			fmt.Printf("Token validation error: %v\n", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		fmt.Printf("Token validated successfully. UserID: %d, Username: %s, Role: %s\n",
			claims.UserID, claims.Username, claims.Role)

		// Добавляем данные пользователя в контекст
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "user_role", claims.Role)

		// Передаем запрос дальше с обновленным контекстом
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
