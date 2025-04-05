package domain

import (
	"context"
	"time"
)

type Article struct {
	ID         string
	Title      string
	UpdatedAt  time.Time
	CreatedAt  time.Time
	Image      string
	Content    string
	LastEditor string
	Creator    string
}

type ArticleInteractor interface {
	CreateArticle(ctx context.Context, title string, image string, content string, creatorName string) (*Article, error)
	Article(ctx context.Context, id string) (*Article, error)
	Articles(ctx context.Context, page, limit int) ([]*Article, error)
	UpdateArticle(ctx context.Context, id string, title string, image string, content string, editorName string) (*Article, error)
	DeteleArticle(ctx context.Context, id string, userID string) error
}

type ArticleRepository interface {
	CreateArticle(ctx context.Context, article *Article) (string, error)
	Article(ctx context.Context, id string) (*Article, error)
	Articles(ctx context.Context, page, limit int) ([]*Article, error)
	UpdateArticle(ctx context.Context, article *Article) (*Article, error)
	DeleteArticle(ctx context.Context, id string) error
}
