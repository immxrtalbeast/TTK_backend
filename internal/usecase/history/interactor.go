package history

import (
	"context"
	"fmt"

	"github.com/immxrtalbeast/TTK_backend/internal/domain"
)

type HistoryInteractor struct {
	historyRepo domain.HistoryRepository
}

func NewHistoryInteractor(historyRepo domain.HistoryRepository) domain.HistoryInteractor {
	return &HistoryInteractor{historyRepo: historyRepo}
}
func (hi *HistoryInteractor) InitHistory(ctx context.Context, articleID string, userID string, articleTitle string) error {
	const op = "uc.history.init"

	history := &domain.History{
		ArticleId:    articleID,
		UserId:       userID,
		ArticleTitle: articleTitle,
	}
	err := hi.historyRepo.InitHistory(ctx, history)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (hi *HistoryInteractor) UpdateHistory(ctx context.Context, articleID string, userID string, eventType domain.EventType, articleTitle string) error {
	const op = "uc.history.update"

	history := &domain.History{
		ArticleId:    articleID,
		UserId:       userID,
		ArticleTitle: articleTitle,
		EventType:    eventType,
	}
	err := hi.historyRepo.UpdateHistory(ctx, history)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

}

func (hi *HistoryInteractor) Histories(ctx context.Context, page, limit int) ([]*domain.History, error) {
	const op = "uc.history.all"
	histories, err := hi.historyRepo.Histories(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return histories, nil
}
