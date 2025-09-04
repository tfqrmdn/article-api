package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"article-api/internal/models"
	"article-api/internal/repository"
)

// MockArticleRepository is a mock implementation of ArticleRepository for testing
type MockArticleRepository struct {
	articles []models.ArticleListItem
	authors  map[string]*models.Author
}

func NewMockArticleRepository() *MockArticleRepository {
	return &MockArticleRepository{
		authors: map[string]*models.Author{
			"author-1": {ID: "author-1", Name: "John Doe"},
			"author-2": {ID: "author-2", Name: "Jane Smith"},
		},
	}
}

func (m *MockArticleRepository) ListArticles(params repository.ListArticlesParams) (*repository.ListArticlesResult, error) {
	// Simple mock implementation - in real tests you'd filter based on params
	return &repository.ListArticlesResult{
		Articles: m.articles,
		Total:    len(m.articles),
		Page:     params.Page,
		Limit:    params.Limit,
	}, nil
}

func (m *MockArticleRepository) CreateArticle(req models.CreateArticleRequest) (*models.Article, error) {
	author, exists := m.authors[req.AuthorID]
	if !exists {
		return nil, &repository.AuthorNotFoundError{}
	}

	article := &models.Article{
		ID:        "test-article-1",
		AuthorID:  req.AuthorID,
		Title:     req.Title,
		Body:      req.Body,
		CreatedAt: time.Now(),
		Author:    author,
	}

	// Convert Article to ArticleListItem for the mock
	articleListItem := models.ArticleListItem{
		ID:        article.ID,
		AuthorID:  article.AuthorID,
		Title:     article.Title,
		CreatedAt: article.CreatedAt,
		Author:    article.Author,
	}
	m.articles = append(m.articles, articleListItem)
	return article, nil
}

func (m *MockArticleRepository) GetAuthorByID(id string) (*models.Author, error) {
	author, exists := m.authors[id]
	if !exists {
		return nil, &repository.AuthorNotFoundError{}
	}
	return author, nil
}

func TestArticleHandler_ListArticles(t *testing.T) {
	mockRepo := NewMockArticleRepository()
	handler := NewArticleHandler(mockRepo)

	// Add some test articles
	mockRepo.articles = []models.ArticleListItem{
		{
			ID:       "article-1",
			AuthorID: "author-1",
			Title:    "Test Article 1",
			Author:   &models.Author{ID: "author-1", Name: "John Doe"},
		},
	}

	req := httptest.NewRequest("GET", "/articles", nil)
	req.Header.Set("X-API-Key", "default-api-key-123")
	w := httptest.NewRecorder()

	handler.ListArticles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var articles []models.ArticleListItem
	if err := json.NewDecoder(w.Body).Decode(&articles); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(articles) != 1 {
		t.Errorf("Expected 1 article, got %d", len(articles))
	}

	if articles[0].Title != "Test Article 1" {
		t.Errorf("Expected title 'Test Article 1', got '%s'", articles[0].Title)
	}
}

func TestArticleHandler_CreateArticle(t *testing.T) {
	mockRepo := NewMockArticleRepository()
	handler := NewArticleHandler(mockRepo)

	reqBody := models.CreateArticleRequest{
		AuthorID: "author-1",
		Title:    "New Test Article",
		Body:     "This is a new test article",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/articles", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "default-api-key-123")
	w := httptest.NewRecorder()

	handler.CreateArticle(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var article models.Article
	if err := json.NewDecoder(w.Body).Decode(&article); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if article.Title != reqBody.Title {
		t.Errorf("Expected title '%s', got '%s'", reqBody.Title, article.Title)
	}

	if article.AuthorID != reqBody.AuthorID {
		t.Errorf("Expected author ID '%s', got '%s'", reqBody.AuthorID, article.AuthorID)
	}
}

func TestArticleHandler_CreateArticle_InvalidJSON(t *testing.T) {
	mockRepo := NewMockArticleRepository()
	handler := NewArticleHandler(mockRepo)

	req := httptest.NewRequest("POST", "/articles", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "default-api-key-123")
	w := httptest.NewRecorder()

	handler.CreateArticle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestArticleHandler_CreateArticle_MissingFields(t *testing.T) {
	mockRepo := NewMockArticleRepository()
	handler := NewArticleHandler(mockRepo)

	reqBody := models.CreateArticleRequest{
		AuthorID: "author-1",
		// Missing title and body
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/articles", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "default-api-key-123")
	w := httptest.NewRecorder()

	handler.CreateArticle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestArticleHandler_CreateArticle_InvalidAuthor(t *testing.T) {
	mockRepo := NewMockArticleRepository()
	handler := NewArticleHandler(mockRepo)

	reqBody := models.CreateArticleRequest{
		AuthorID: "non-existent-author",
		Title:    "Test Article",
		Body:     "Test body",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/articles", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "default-api-key-123")
	w := httptest.NewRecorder()

	handler.CreateArticle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}
