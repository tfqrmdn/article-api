package repository

import (
	"database/sql"
	"testing"

	"article-api/internal/cache"
	"article-api/internal/models"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Connect to test database
	db, err := sql.Open("postgres", "host=localhost port=5432 user=article_user password=article_password dbname=article_db sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean up test data
	_, err = db.Exec("DELETE FROM articles WHERE id LIKE 'test-%'")
	if err != nil {
		t.Fatalf("Failed to clean up test data: %v", err)
	}

	return db
}

func TestArticleRepository_ListArticles(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	mockCache := cache.NewMockCacheService()
	repo := NewArticleRepository(db, mockCache)

	// Test listing articles with default parameters
	params := ListArticlesParams{
		Page:  1,
		Limit: 10,
	}
	result, err := repo.ListArticles(params)
	if err != nil {
		t.Fatalf("Failed to list articles: %v", err)
	}

	// Should have at least the sample data
	if len(result.Articles) < 2 {
		t.Errorf("Expected at least 2 articles, got %d", len(result.Articles))
	}

	// Check that articles have author information
	for _, article := range result.Articles {
		if article.Author == nil {
			t.Errorf("Article %s should have author information", article.ID)
		}
		if article.Author.ID == "" {
			t.Errorf("Article %s author should have ID", article.ID)
		}
		if article.Author.Name == "" {
			t.Errorf("Article %s author should have name", article.ID)
		}
	}

	// Test pagination
	params.Limit = 1
	result, err = repo.ListArticles(params)
	if err != nil {
		t.Fatalf("Failed to list articles with pagination: %v", err)
	}
	if len(result.Articles) != 1 {
		t.Errorf("Expected 1 article with limit=1, got %d", len(result.Articles))
	}
}

func TestArticleRepository_CreateArticle(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	mockCache := cache.NewMockCacheService()
	repo := NewArticleRepository(db, mockCache)

	req := models.CreateArticleRequest{
		AuthorID: "author-1",
		Title:    "Test Article",
		Body:     "This is a test article body",
	}

	article, err := repo.CreateArticle(req)
	if err != nil {
		t.Fatalf("Failed to create article: %v", err)
	}

	if article.ID == "" {
		t.Error("Created article should have an ID")
	}
	if article.AuthorID != req.AuthorID {
		t.Errorf("Expected author ID %s, got %s", req.AuthorID, article.AuthorID)
	}
	if article.Title != req.Title {
		t.Errorf("Expected title %s, got %s", req.Title, article.Title)
	}
	if article.Body != req.Body {
		t.Errorf("Expected body %s, got %s", req.Body, article.Body)
	}
	if article.Author == nil {
		t.Error("Created article should have author information")
	}
	if article.Author.ID != req.AuthorID {
		t.Errorf("Expected author ID %s, got %s", req.AuthorID, article.Author.ID)
	}
}

func TestArticleRepository_CreateArticle_InvalidAuthor(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	mockCache := cache.NewMockCacheService()
	repo := NewArticleRepository(db, mockCache)

	req := models.CreateArticleRequest{
		AuthorID: "non-existent-author",
		Title:    "Test Article",
		Body:     "This is a test article body",
	}

	_, err := repo.CreateArticle(req)
	if err == nil {
		t.Error("Expected error when creating article with non-existent author")
	}
}

func TestArticleRepository_GetAuthorByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	mockCache := cache.NewMockCacheService()
	repo := NewArticleRepository(db, mockCache)

	// Test existing author
	author, err := repo.GetAuthorByID("author-1")
	if err != nil {
		t.Fatalf("Failed to get author: %v", err)
	}

	if author.ID != "author-1" {
		t.Errorf("Expected author ID author-1, got %s", author.ID)
	}
	if author.Name == "" {
		t.Error("Author should have a name")
	}

	// Test non-existent author
	_, err = repo.GetAuthorByID("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent author")
	}
}
