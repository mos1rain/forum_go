package service

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrEmptyContent  = errors.New("message content cannot be empty")
	ErrEmptyUsername = errors.New("username cannot be empty")
	ErrInvalidUserID = errors.New("invalid user ID")
)

type Message struct {
	ID        int
	UserID    int
	Username  string
	Content   string
	CreatedAt time.Time
}

type ChatService struct {
	db *sql.DB
}

func NewChatService(db *sql.DB) *ChatService {
	return &ChatService{
		db: db,
	}
}

func (c *ChatService) AddMessage(userID int, username, content string) (Message, error) {
	if content == "" {
		return Message{}, ErrEmptyContent
	}
	if username == "" {
		return Message{}, ErrEmptyUsername
	}
	if userID <= 0 {
		return Message{}, ErrInvalidUserID
	}

	var msg Message
	query := `
		INSERT INTO chat_messages (user_id, username, content)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, username, content, created_at
	`
	err := c.db.QueryRow(query, userID, username, content).Scan(
		&msg.ID,
		&msg.UserID,
		&msg.Username,
		&msg.Content,
		&msg.CreatedAt,
	)
	return msg, err
}

func (c *ChatService) GetHistory(limit int) ([]Message, error) {
	if limit < 0 {
		limit = 0
	}

	query := `
		SELECT id, user_id, username, content, created_at
		FROM chat_messages
		ORDER BY created_at ASC
		LIMIT $1
	`
	rows, err := c.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(
			&msg.ID,
			&msg.UserID,
			&msg.Username,
			&msg.Content,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (c *ChatService) DeleteMessage(id int) error {
	if id <= 0 {
		return ErrInvalidUserID
	}
	query := `DELETE FROM chat_messages WHERE id = $1`
	_, err := c.db.Exec(query, id)
	return err
}

func (c *ChatService) CleanOldMessages(olderThan time.Duration) (int, error) {
	query := `
		DELETE FROM chat_messages
		WHERE created_at < NOW() - INTERVAL '1 day'
		RETURNING id
	`
	result, err := c.db.Exec(query)
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	return int(count), err
}
