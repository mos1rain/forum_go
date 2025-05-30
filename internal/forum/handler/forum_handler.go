package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mos1rain/forum_go/internal/forum/models"
	"github.com/mos1rain/forum_go/internal/forum/service"
)

type ForumHandler struct {
	service *service.ForumService
}

func NewForumHandler(service *service.ForumService) *ForumHandler {
	return &ForumHandler{service: service}
}

func (h *ForumHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.Categories.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *ForumHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем ID пользователя из контекста
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	category.CreatorID = int64(userID)
	if err := h.service.Categories.Create(r.Context(), &category); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *ForumHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.Posts.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *ForumHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Сначала декодируем запрос в map для обработки строковых значений
	var requestData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем ID пользователя из контекста
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Создаем объект Post
	var post models.Post
	post.Title = requestData["title"].(string)
	post.Content = requestData["content"].(string)

	// Преобразуем categoryId из строки в число
	categoryIDStr, ok := requestData["categoryId"].(string)
	if !ok {
		http.Error(w, "Invalid category ID format", http.StatusBadRequest)
		return
	}
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	post.CategoryID = categoryID

	// Устанавливаем ID пользователя
	post.AuthorID = int64(userID)

	if err := h.service.Posts.Create(&post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *ForumHandler) GetPostByID(id int) (*models.Post, error) {
	return h.service.Posts.GetByID(id)
}

func (h *ForumHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	if err := h.service.Posts.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
}

func (h *ForumHandler) GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	comments, err := h.service.Comments.GetByPostID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (h *ForumHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.Comments.Create(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *ForumHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	if err := h.service.Comments.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment deleted successfully"})
}

func (h *ForumHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Получаем роль пользователя из контекста
	userRole, ok := r.Context().Value("user_role").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.Categories.Delete(r.Context(), id, userRole); err != nil {
		switch {
		case errors.Is(err, service.ErrAdminRoleRequired):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, service.ErrCategoryNotFound):
			http.Error(w, "Category not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully"})
}
