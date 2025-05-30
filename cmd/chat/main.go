package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	_ "github.com/mos1rain/forum_go/docs"
	"github.com/mos1rain/forum_go/internal/chat/service"
	"github.com/rs/zerolog"
	_ "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "modernc.org/sqlite"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan service.Message)
	mutex     sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Подключение к SQLite
	db, err := sql.Open("sqlite", "./forum.db")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to ping database")
	}

	// Инициализация таблицы chat_messages
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chat_messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			username TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize chat_messages table")
	}

	chatService := service.NewChatService(db)

	go handleMessages(chatService)
	go cleanOldMessages(chatService)

	http.HandleFunc("/ws", withCORS(handleWS(chatService)))

	http.HandleFunc("/history", withCORS(func(w http.ResponseWriter, r *http.Request) {
		history, err := chatService.GetHistory(50)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get chat history")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		out := make([]map[string]interface{}, 0, len(history))
		for _, msg := range history {
			out = append(out, map[string]interface{}{
				"id":         msg.ID,
				"user_id":    msg.UserID,
				"username":   msg.Username,
				"content":    msg.Content,
				"created_at": msg.CreatedAt,
			})
		}
		json.NewEncoder(w).Encode(out)
	}))

	http.HandleFunc("/messages", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var m service.Message
			if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if m.Content == "" || m.Username == "" || m.UserID == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			msg, err := chatService.AddMessage(m.UserID, m.Username, m.Content)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add message")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			broadcast <- msg
			w.WriteHeader(http.StatusCreated)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))

	http.HandleFunc("/delete_message", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := parseJWT(token)
		if err != nil || claims["role"] != "admin" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := chatService.DeleteMessage(id); err != nil {
			logger.Error().Err(err).Msg("Failed to delete message")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	logger.Info().Msg("Chat service started on :3003")
	logger.Fatal().Err(http.ListenAndServe(":3003", nil)).Msg("chat server crashed")
}

func handleWS(chatService *service.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error().Err(err).Msg("WebSocket upgrade error")
			return
		}
		defer func() {
			conn.Close()
			delete(clients, conn)
		}()

		clients[conn] = true

		// Отправляем историю сообщений при подключении
		history, err := chatService.GetHistory(50)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get chat history")
		} else {
			for _, msg := range history {
				logger.Info().Msgf("Send history to client: %+v", msg)
				out := map[string]interface{}{
					"id":         msg.ID,
					"user_id":    msg.UserID,
					"username":   msg.Username,
					"content":    msg.Content,
					"created_at": msg.CreatedAt,
				}
				data, err := json.Marshal(out)
				if err != nil {
					logger.Error().Err(err).Msg("marshal error")
					continue
				}
				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					logger.Error().Err(err).Msg("Failed to send history")
					return
				}
			}
		}

		for {
			var raw map[string]interface{}
			if err := conn.ReadJSON(&raw); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Warn().Err(err).Msg("WebSocket read error")
				}
				return
			}

			logger.Info().Msgf("RAW from client: %+v", raw)
			userID := 0
			if v, ok := raw["user_id"].(float64); ok {
				userID = int(v)
			}
			username, _ := raw["username"].(string)
			content, _ := raw["content"].(string)

			if content == "" || username == "" {
				logger.Warn().Msg("Empty message content or username")
				continue
			}

			msg, err := chatService.AddMessage(userID, username, content)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add message from WebSocket")
				continue
			}
			logger.Info().Msgf("Broadcast to clients: %+v", msg)
			broadcast <- msg
		}
	}
}

func handleMessages(chatService *service.ChatService) {
	for {
		msg := <-broadcast
		logger.Info().Msgf("Send to client: %+v", msg)
		mutex.Lock()
		for client := range clients {
			out := map[string]interface{}{
				"id":         msg.ID,
				"user_id":    msg.UserID,
				"username":   msg.Username,
				"content":    msg.Content,
				"created_at": msg.CreatedAt,
			}
			logger.Info().Msgf("Send to client (WS): %+v", out)
			data, err := json.Marshal(out)
			if err != nil {
				logger.Error().Err(err).Msg("marshal error")
				continue
			}
			if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
				logger.Warn().Err(err).Msg("WebSocket write error, disconnecting client")
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func cleanOldMessages(chatService *service.ChatService) {
	for {
		removed, err := chatService.CleanOldMessages(24 * time.Hour)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to clean old messages")
		} else if removed > 0 {
			logger.Info().Int("count", removed).Msg("Cleaned old messages")
		}
		time.Sleep(5 * time.Minute)
	}
}

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

func parseJWT(token string) (map[string]interface{}, error) {
	parsed, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
		out := make(map[string]interface{})
		for k, v := range claims {
			out[k] = v
		}
		return out, nil
	}
	return nil, errors.New("invalid token claims")
}
