package repository

import "article-api/internal/models"

// ArticleRepositoryInterface defines the contract for article repository operations
type ArticleRepositoryInterface interface {
	ListArticles(params ListArticlesParams) (*ListArticlesResult, error)
	CreateArticle(req models.CreateArticleRequest) (*models.Article, error)
	GetAuthorByID(id string) (*models.Author, error)
}

// ListArticlesParams holds parameters for listing articles
type ListArticlesParams struct {
	Search     string
	AuthorName string
	Page       int
	Limit      int
}

// ListArticlesResult holds the result of listing articles
type ListArticlesResult struct {
	Articles []models.ArticleListItem
	Total    int
	Page     int
	Limit    int
}

// AuthorNotFoundError represents an error when author is not found
type AuthorNotFoundError struct{}

func (e *AuthorNotFoundError) Error() string {
	return "author not found"
}
