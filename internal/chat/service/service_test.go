package service

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:28072005@localhost:5432/forum_test?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean up and create tables
	_, err = db.Exec(`
		DROP TABLE IF EXISTS chat_messages;
		CREATE TABLE chat_messages (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			username VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	return db
}

func TestAddAndGetHistory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	cs := NewChatService(db)

	// Add messages
	_, err := cs.AddMessage(1, "user1", "hello")
	if err != nil {
		t.Fatalf("Failed to add message: %v", err)
	}

	_, err = cs.AddMessage(2, "user2", "world")
	if err != nil {
		t.Fatalf("Failed to add message: %v", err)
	}

	// Get history
	msgs, err := cs.GetHistory(2)
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}

	if len(msgs) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(msgs))
	}

	// Проверяем, что оба сообщения есть, независимо от порядка
	foundHello := false
	foundWorld := false
	for _, m := range msgs {
		if m.Content == "hello" {
			foundHello = true
		}
		if m.Content == "world" {
			foundWorld = true
		}
	}
	if !foundHello || !foundWorld {
		t.Errorf("Expected both 'hello' and 'world' messages, got: %+v", msgs)
	}
}

func TestDeleteMessage(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	cs := NewChatService(db)

	// Add message
	msg, err := cs.AddMessage(1, "user1", "to delete")
	if err != nil {
		t.Fatalf("Failed to add message: %v", err)
	}

	// Delete message
	err = cs.DeleteMessage(msg.ID)
	if err != nil {
		t.Fatalf("Failed to delete message: %v", err)
	}

	// Check history
	msgs, err := cs.GetHistory(10)
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}

	if len(msgs) != 0 {
		t.Errorf("Expected 0 messages, got %d", len(msgs))
	}
}

func TestAddMessageValidation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	cs := NewChatService(db)

	// Test empty content
	_, err := cs.AddMessage(1, "user1", "")
	if err == nil {
		t.Error("Expected error for empty content, got nil")
	}

	// Test empty username
	_, err = cs.AddMessage(1, "", "test message")
	if err == nil {
		t.Error("Expected error for empty username, got nil")
	}

	// Test invalid user ID
	_, err = cs.AddMessage(-1, "user1", "test message")
	if err == nil {
		t.Error("Expected error for negative user ID, got nil")
	}
}

func TestGetHistoryEdgeCases(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	cs := NewChatService(db)

	// Test with zero limit
	msgs, err := cs.GetHistory(0)
	if err != nil {
		t.Fatalf("Failed to get history with zero limit: %v", err)
	}
	if len(msgs) != 0 {
		t.Errorf("Expected 0 messages with zero limit, got %d", len(msgs))
	}

	// Test with negative limit
	msgs, err = cs.GetHistory(-1)
	if err != nil {
		t.Fatalf("Failed to get history with negative limit: %v", err)
	}
	if len(msgs) != 0 {
		t.Errorf("Expected 0 messages with negative limit, got %d", len(msgs))
	}

	// Add some messages
	_, err = cs.AddMessage(1, "user1", "message1")
	if err != nil {
		t.Fatalf("Failed to add message: %v", err)
	}
	_, err = cs.AddMessage(2, "user2", "message2")
	if err != nil {
		t.Fatalf("Failed to add message: %v", err)
	}

	// Test with limit larger than available messages
	msgs, err = cs.GetHistory(10)
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}
	if len(msgs) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(msgs))
	}
}

func TestDeleteMessageEdgeCases(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	cs := NewChatService(db)

	// Test deleting non-existent message
	err := cs.DeleteMessage(999)
	if err != nil {
		t.Errorf("Expected no error when deleting non-existent message, got %v", err)
	}

	// Test deleting with invalid ID
	err = cs.DeleteMessage(-1)
	if err == nil {
		t.Error("Expected error when deleting with negative ID, got nil")
	}

	// Add and then delete a message
	msg, err := cs.AddMessage(1, "user1", "to delete")
	if err != nil {
		t.Fatalf("Failed to add message: %v", err)
	}

	err = cs.DeleteMessage(msg.ID)
	if err != nil {
		t.Fatalf("Failed to delete message: %v", err)
	}

	// Verify message is deleted
	msgs, err := cs.GetHistory(10)
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}
	if len(msgs) != 0 {
		t.Errorf("Expected 0 messages after deletion, got %d", len(msgs))
	}
}

func TestCleanOldMessages(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	cs := NewChatService(db)

	// Add a message
	_, err := cs.AddMessage(1, "user1", "test message")
	if err != nil {
		t.Fatalf("Failed to add message: %v", err)
	}

	// Clean messages
	removed, err := cs.CleanOldMessages(24 * time.Hour)
	if err != nil {
		t.Fatalf("Failed to clean old messages: %v", err)
	}

	// Since the message is new, it shouldn't be removed
	if removed != 0 {
		t.Errorf("Expected 0 messages removed, got %d", removed)
	}

	// Verify message still exists
	msgs, err := cs.GetHistory(10)
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}
	if len(msgs) != 1 {
		t.Errorf("Expected 1 message to remain, got %d", len(msgs))
	}
}
