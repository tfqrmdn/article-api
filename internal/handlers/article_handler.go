package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"article-api/internal/models"
	"article-api/internal/repository"
)

// ArticleHandler handles HTTP requests for articles
type ArticleHandler struct {
	repo repository.ArticleRepositoryInterface
}

// NewArticleHandler creates a new article handler
func NewArticleHandler(repo repository.ArticleRepositoryInterface) *ArticleHandler {
	return &ArticleHandler{repo: repo}
}

// ListArticles handles GET /articles
func (h *ArticleHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	search := r.URL.Query().Get("search")
	authorName := r.URL.Query().Get("author")
	page := parseIntParam(r.URL.Query().Get("page"), 1)
	limit := parseIntParam(r.URL.Query().Get("limit"), 10)

	// Create parameters
	params := repository.ListArticlesParams{
		Search:     search,
		AuthorName: authorName,
		Page:       page,
		Limit:      limit,
	}

	result, err := h.repo.ListArticles(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list articles: %v", err), http.StatusInternalServerError)
		return
	}

	// Set pagination headers
	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", result.Total))
	w.Header().Set("X-Page", fmt.Sprintf("%d", result.Page))
	w.Header().Set("X-Limit", fmt.Sprintf("%d", result.Limit))
	w.Header().Set("X-Total-Pages", fmt.Sprintf("%d", (result.Total+result.Limit-1)/result.Limit))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result.Articles); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// parseIntParam parses an integer parameter with a default value
func parseIntParam(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
		return parsed
	}
	return defaultValue
}

// CreateArticle handles POST /articles
func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var req models.CreateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.AuthorID == "" || req.Title == "" || req.Body == "" {
		http.Error(w, "Missing required fields: author_id, title, body", http.StatusBadRequest)
		return
	}

	// Check if author exists
	_, err := h.repo.GetAuthorByID(req.AuthorID)
	if err != nil {
		http.Error(w, "Author not found", http.StatusBadRequest)
		return
	}

	article, err := h.repo.CreateArticle(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create article: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(article); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
