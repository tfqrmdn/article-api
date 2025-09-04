package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"article-api/internal/cache"
	"article-api/internal/models"
)

// ArticleRepository handles database operations for articles
type ArticleRepository struct {
	db    *sql.DB
	cache cache.CacheServiceInterface
}

// NewArticleRepository creates a new article repository
func NewArticleRepository(db *sql.DB, cacheService cache.CacheServiceInterface) *ArticleRepository {
	return &ArticleRepository{
		db:    db,
		cache: cacheService,
	}
}

// ListArticles retrieves articles with search, filtering, and pagination
func (r *ArticleRepository) ListArticles(params ListArticlesParams) (*ListArticlesResult, error) {
	// Set defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100 // Max limit
	}

	offset := (params.Page - 1) * params.Limit

	// Build WHERE clause
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if params.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(a.title ILIKE $%d OR a.body ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	if params.AuthorName != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("au.name ILIKE $%d", argIndex))
		args = append(args, "%"+params.AuthorName+"%")
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM articles a
		LEFT JOIN authors au ON a.author_id = au.id
		%s
	`, whereClause)

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count articles: %w", err)
	}

	// Articles query with pagination (excluding body for performance)
	articlesQuery := fmt.Sprintf(`
		SELECT 
			a.id, 
			a.author_id, 
			a.title, 
			a.created_at,
			au.id as author_id,
			au.name as author_name
		FROM articles a
		LEFT JOIN authors au ON a.author_id = au.id
		%s
		ORDER BY a.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(articlesQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query articles: %w", err)
	}
	defer rows.Close()

	var articles []models.ArticleListItem
	for rows.Next() {
		var article models.ArticleListItem
		var author models.Author

		err := rows.Scan(
			&article.ID,
			&article.AuthorID,
			&article.Title,
			&article.CreatedAt,
			&author.ID,
			&author.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)
		}

		article.Author = &author
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating articles: %w", err)
	}

	return &ListArticlesResult{
		Articles: articles,
		Total:    total,
		Page:     params.Page,
		Limit:    params.Limit,
	}, nil
}

// CreateArticle creates a new article
func (r *ArticleRepository) CreateArticle(req models.CreateArticleRequest) (*models.Article, error) {
	// Generate a simple ID (in production, you might want to use UUID)
	id := fmt.Sprintf("article-%d", time.Now().UnixNano())

	query := `
		INSERT INTO articles (id, author_id, title, body, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, author_id, title, body, created_at
	`

	var article models.Article
	err := r.db.QueryRow(query, id, req.AuthorID, req.Title, req.Body, time.Now()).
		Scan(&article.ID, &article.AuthorID, &article.Title, &article.Body, &article.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	// Fetch author information
	authorQuery := `SELECT id, name FROM authors WHERE id = $1`
	var author models.Author
	err = r.db.QueryRow(authorQuery, req.AuthorID).
		Scan(&author.ID, &author.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch author: %w", err)
	}

	article.Author = &author

	// Cache the created article for 10 minutes (600 seconds)
	cacheKey := fmt.Sprintf("article:%s", article.ID)
	if cacheErr := r.cache.SetWithTTL(cacheKey, article, 600); cacheErr != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to cache created article: %v\n", cacheErr)
	}

	// Invalidate cache when new article is created
	if cacheErr := r.cache.Delete("articles:list"); cacheErr != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to invalidate cache: %v\n", cacheErr)
	}

	return &article, nil
}

// GetAuthorByID retrieves an author by ID
func (r *ArticleRepository) GetAuthorByID(id string) (*models.Author, error) {
	query := `SELECT id, name FROM authors WHERE id = $1`

	var author models.Author
	err := r.db.QueryRow(query, id).
		Scan(&author.ID, &author.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &AuthorNotFoundError{}
		}
		return nil, fmt.Errorf("failed to get author: %w", err)
	}

	return &author, nil
}
