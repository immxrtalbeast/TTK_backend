package domain

import (
	"context"
	"time"
)

type Article struct {
	ID               string
	Title            string
	UpdatedAt        time.Time
	CreatedAt        time.Time
	Name_last_reedit string
	Image            string
	Content          string
}

type ArticleInteractor interface {
	CreateArticle(ctx context.Context, title string, image string, content string, name_last_edit string) (*Article, error)
	Article(ctx context.Context, id string) (*Article, error)
	UpdateArticle(ctx context.Context, title string, image string, content string, name_last_edit string) (*Article, error)
}

type ArticleRepository interface {
	CreateArticle(ctx context.Context, article *Article) (string, error)
	Article(ctx context.Context, id string) (*Article, error)
	// UpdateArticleContent(ctx context.Context, name_last_reedit string, content string) error
	// DeleteArticle(ctx context.Context, id string) error
}
