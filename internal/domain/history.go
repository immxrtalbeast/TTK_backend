package domain

import (
	"context"
	"time"
)

type EventType string

const (
	Changed EventType = "CHANGED"
	Delete  EventType = "DELETE"
	Create  EventType = "CREATE"
)

type History struct {
	ID           string
	ArticleId    string
	UserId       string
	ChangedAt    time.Time
	EventType    EventType
	ArticleTitle string
}

type HistoryInteractor interface {
	InitHistory(ctx context.Context, articleID string, userID string, articleTitle string) error
	UpdateHistory(ctx context.Context, articleID string, userID string, eventType EventType, articleTitle string) error
	Histories(ctx context.Context, page, limit int) ([]*History, error)
}

type HistoryRepository interface {
	InitHistory(ctx context.Context, history *History) error
	UpdateHistory(ctx context.Context, history *History) error
	Histories(ctx context.Context, page, limit int) ([]*History, error)
}
