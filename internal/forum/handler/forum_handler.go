package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mos1rain/forum_go/internal/forum/middleware"
	"github.com/mos1rain/forum_go/internal/forum/models"
	"github.com/mos1rain/forum_go/internal/forum/service"
	"github.com/mos1rain/forum_go/pkg/jwt"
)

type ForumHandler struct {
	service *service.ForumService
}

func NewForumHandler(service *service.ForumService) *ForumHandler {
	return &ForumHandler{service: service}
}

// --- Категории ---
// @Summary Создать категорию
// @Description Создать новую категорию форума
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Категория"
// @Success 201 {object} models.Category
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/forum/categories [post]
func (h *ForumHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input models.Category
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.service.Categories.Create(&input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

// @Summary Получить список категорий
// @Description Получить все категории форума
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category
// @Failure 500 {string} string
// @Router /api/forum/categories [get]
func (h *ForumHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	categories, err := h.service.Categories.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if categories == nil {
		categories = []models.Category{}
	}
	json.NewEncoder(w).Encode(categories)
}

// --- Посты ---
// @Summary Создать пост
// @Description Создать новый пост в категории
// @Tags posts
// @Accept json
// @Produce json
// @Param post body models.Post true "Пост"
// @Success 201 {object} models.Post
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /api/forum/posts [post]
func (h *ForumHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input models.Post
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	input.UserID = userID
	if err := h.service.Posts.Create(&input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

// @Summary Получить список постов
// @Description Получить все посты
// @Tags posts
// @Produce json
// @Success 200 {array} models.Post
// @Failure 500 {string} string
// @Router /api/forum/posts [get]
func (h *ForumHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	posts, err := h.service.Posts.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if posts == nil {
		posts = []models.Post{}
	}
	json.NewEncoder(w).Encode(posts)
}

// --- Комментарии ---
// @Summary Создать комментарий
// @Description Добавить комментарий к посту
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body models.Comment true "Комментарий"
// @Success 201 {object} models.Comment
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /api/forum/comments [post]
func (h *ForumHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input models.Comment
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	input.UserID = userID
	if err := h.service.Comments.Create(&input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

// @Summary Получить комментарии к посту
// @Description Получить все комментарии по post_id
// @Tags comments
// @Produce json
// @Param post_id query int true "ID поста"
// @Success 200 {array} models.Comment
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/forum/comments [get]
func (h *ForumHandler) GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}
	comments, err := h.service.Comments.GetByPostID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if comments == nil {
		comments = []models.Comment{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (h *ForumHandler) GetPostByID(id int) (*models.Post, error) {
	return h.service.Posts.GetByID(id)
}

// @Summary Удалить пост
// @Description Удалить пост по id (только для админа)
// @Tags posts
// @Produce json
// @Param id query int true "ID поста"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /api/forum/delete_post [delete]
func (h *ForumHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims, ok := r.Context().Value("claims").(*jwt.Claims)
	if !ok || claims.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	if err := h.service.Posts.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Удалить комментарий
// @Description Удалить комментарий по id (только для админа)
// @Tags comments
// @Produce json
// @Param id query int true "ID комментария"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /api/forum/delete_comment [delete]
func (h *ForumHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims, ok := r.Context().Value("claims").(*jwt.Claims)
	if !ok || claims.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	if err := h.service.Comments.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Удалить категорию
// @Description Удалить категорию по id (только для админа)
// @Tags categories
// @Produce json
// @Param id query int true "ID категории"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /api/forum/delete_category [delete]
func (h *ForumHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims, ok := r.Context().Value("claims").(*jwt.Claims)
	if !ok || claims.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	if err := h.service.Categories.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
