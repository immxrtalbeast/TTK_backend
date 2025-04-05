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

func (ai *ArticleInteractor) CreateArticle(ctx context.Context, title string, image string, content string, name_last_edit string) (*domain.Article, error) {
	const op = "uc.article.create"
	article := domain.Article{
		Title:            title,
		Image:            image,
		Content:          content,
		Name_last_reedit: name_last_edit,
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

func (ai *ArticleInteractor) UpdateArticle(ctx context.Context, title string, image string, content string, name_last_edit string) (*domain.Article, error) {
	const op = "uc.article.update"
	return nil, nil
}
