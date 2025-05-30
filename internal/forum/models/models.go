package models

import "time"

// Category represents a forum category
// @Description Forum category information
type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`        // Название категории
	Description string    `json:"description"` // Описание категории
	CreatorID   int64     `json:"creator_id"`  // ID создателя категории
	CreatedAt   time.Time `json:"created_at"`  // Дата создания
	UpdatedAt   time.Time `json:"updated_at"`  // Дата последнего обновления
}

// Post represents a forum post
// @Description Forum post information
type Post struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`       // Заголовок поста
	Content    string    `json:"content"`     // Содержание поста
	CategoryID int64     `json:"category_id"` // ID категории
	AuthorID   int64     `json:"author_id"`   // ID автора
	CreatedAt  time.Time `json:"created_at"`  // Дата создания
	UpdatedAt  time.Time `json:"updated_at"`  // Дата последнего обновления
}

// Comment represents a forum comment
// @Description Forum comment information
type Comment struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`    // Содержание комментария
	PostID    int64     `json:"post_id"`    // ID поста
	AuthorID  int64     `json:"author_id"`  // ID автора
	CreatedAt time.Time `json:"created_at"` // Дата создания
	UpdatedAt time.Time `json:"updated_at"` // Дата последнего обновления
}
