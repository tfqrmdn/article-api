package models

import "time"

// Author represents an author in the system
type Author struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Article represents an article in the system
type Article struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	Author    *Author   `json:"author,omitempty"`
}

// ArticleListItem represents an article in list responses (without body for performance)
type ArticleListItem struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Author    *Author   `json:"author,omitempty"`
}

// CreateArticleRequest represents the request payload for creating an article
type CreateArticleRequest struct {
	AuthorID string `json:"author_id" validate:"required"`
	Title    string `json:"title" validate:"required"`
	Body     string `json:"body" validate:"required"`
}
