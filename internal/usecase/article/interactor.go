package article

import (
	"context"
	"errors"
	"fmt"

	"github.com/immxrtalbeast/TTK_backend/internal/domain"
	"github.com/immxrtalbeast/TTK_backend/storage/prisma/db"
)

type ArticleInteractor struct {
	articleRepo domain.ArticleRepository
}

func NewArticleInteractor(articleRepo domain.ArticleRepository) domain.ArticleInteractor {
	return &ArticleInteractor{articleRepo: articleRepo}
}

func (ai *ArticleInteractor) Article(ctx context.Context, id string) (*domain.Article, error) {
	const op = "uc.article.get"

	article, err := ai.articleRepo.Article(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, db.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)

	}
	return article, nil
}
func (ai *ArticleInteractor) Articles(ctx context.Context, page, limit int) ([]*domain.Article, error) {
	const op = "uc.article.get.all"
	articles, err := ai.articleRepo.Articles(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return articles, nil
}

func (ai *ArticleInteractor) CreateArticle(ctx context.Context, title string, image string, content string, creatorName string) (*domain.Article, error) {
	const op = "uc.article.create"
	article := domain.Article{
		Title:      title,
		Image:      image,
		Content:    content,
		Creator:    creatorName,
		LastEditor: creatorName,
	}
	articleID, err := ai.articleRepo.CreateArticle(ctx, &article)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}
	articleDB, err := ai.articleRepo.Article(ctx, articleID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return articleDB, nil
}

func (ai *ArticleInteractor) UpdateArticle(ctx context.Context, id string, title string, image string, content string, editorName string) (*domain.Article, error) {
	const op = "uc.article.update"
	article := domain.Article{
		ID:         id,
		Title:      title,
		Image:      image,
		Content:    content,
		LastEditor: editorName,
	}
	result, err := ai.articleRepo.UpdateArticle(ctx, &article)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return result, nil
}
func (ai *ArticleInteractor) DeteleArticle(ctx context.Context, id string, userID string) error {
	const op = "uc.article.delete"
	articleDB, err := ai.articleRepo.Article(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if articleDB.Creator == userID {
		err := ai.articleRepo.DeleteArticle(ctx, id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		return fmt.Errorf("No rights to delete article.")
	}
	return nil
}
