package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/mos1rain/forum_go/internal/auth/models"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Создаем таблицу users
	_, err = db.Exec(`
		CREATE TABLE users (
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
		t.Fatalf("Failed to create users table: %v", err)
	}

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &models.User{
				Username:     "testuser",
				Email:        "test@example.com",
				PasswordHash: "hashedpassword",
				Role:         "user",
			},
			wantErr: false,
		},
		{
			name: "duplicate username",
			user: &models.User{
				Username:     "testuser",
				Email:        "another@example.com",
				PasswordHash: "hashedpassword",
				Role:         "user",
			},
			wantErr: true,
		},
		{
			name: "duplicate email",
			user: &models.User{
				Username:     "anotheruser",
				Email:        "test@example.com",
				PasswordHash: "hashedpassword",
				Role:         "user",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if tt.user.ID == 0 {
					t.Error("Create() didn't set user ID")
				}
				if tt.user.CreatedAt.IsZero() {
					t.Error("Create() didn't set CreatedAt")
				}
				if tt.user.UpdatedAt.IsZero() {
					t.Error("Create() didn't set UpdatedAt")
				}
			}
		})
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Создаем тестового пользователя
	testUser := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := repo.Create(testUser); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name     string
		username string
		want     *models.User
		wantErr  bool
	}{
		{
			name:     "existing user",
			username: "testuser",
			want:     testUser,
			wantErr:  false,
		},
		{
			name:     "non-existing user",
			username: "nonexistent",
			want:     nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil && got != nil {
				t.Errorf("GetByUsername() got = %v, want nil", got)
			}
			if tt.want != nil && got != nil {
				if got.Username != tt.want.Username {
					t.Errorf("GetByUsername() got username = %v, want %v", got.Username, tt.want.Username)
				}
				if got.Email != tt.want.Email {
					t.Errorf("GetByUsername() got email = %v, want %v", got.Email, tt.want.Email)
				}
			}
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Создаем тестового пользователя
	testUser := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         "user",
	}
	if err := repo.Create(testUser); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		email   string
		want    *models.User
		wantErr bool
	}{
		{
			name:    "existing email",
			email:   "test@example.com",
			want:    testUser,
			wantErr: false,
		},
		{
			name:    "non-existing email",
			email:   "nonexistent@example.com",
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil && got != nil {
				t.Errorf("GetByEmail() got = %v, want nil", got)
			}
			if tt.want != nil && got != nil {
				if got.Email != tt.want.Email {
					t.Errorf("GetByEmail() got email = %v, want %v", got.Email, tt.want.Email)
				}
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Создаем тестового пользователя
	testUser := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         "user",
	}
	if err := repo.Create(testUser); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		want    *models.User
		wantErr bool
	}{
		{
			name:    "existing id",
			id:      testUser.ID,
			want:    testUser,
			wantErr: false,
		},
		{
			name:    "non-existing id",
			id:      9999,
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil && got != nil {
				t.Errorf("GetByID() got = %v, want nil", got)
			}
			if tt.want != nil && got != nil {
				if got.ID != tt.want.ID {
					t.Errorf("GetByID() got ID = %v, want %v", got.ID, tt.want.ID)
				}
			}
		})
	}
}
