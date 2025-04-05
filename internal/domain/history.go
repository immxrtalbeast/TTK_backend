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

type HistoryRepository interface {
	InitHistory(ctx context.Context, article *Article) error
	HistoryByID(ctx context.Context, id string) (*History, error)
	UpdateHistory(ctx context.Context, articleID string, userID string, eventType EventType, articleTitle string) error
}
