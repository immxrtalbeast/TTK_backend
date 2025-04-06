package domain

import (
	"context"
	"time"
)

type Status string

const (
	Current   Status = "CURRENT"
	Pending   Status = "PENDING"
	Completed Status = "COMPLETED"
)

type Priority string

const (
	High   Priority = "HIGH"
	Middle Priority = "MIDDLE"
	Low    Priority = "LOW"
)

type Task struct {
	ID        string
	Title     string
	Content   string
	Image     string
	CreatedAt time.Time
	UserID    string
	PlannedAt time.Time
	Priority  Priority
	Status    Status
}

type TaskInteractor interface {
	CreateTask(ctx context.Context, title string, image string, content string, userID string, priority Priority, status Status) (string, error)
	Task(ctx context.Context, id string) (*Task, error)
	Tasks(ctx context.Context, page, limit int) ([]*Task, error)
	UpdateTask(ctx context.Context, id string, title string, image string, content string, userID string, priority Priority, status Status) error
	DeleteTask(ctx context.Context, id string) error
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task *Task) (string, error)
	Task(ctx context.Context, id string) (*Task, error)
	Tasks(ctx context.Context, page, limit int) ([]*Task, error)
	UpdateTask(ctx context.Context, task *Task) error
	DeleteTask(ctx context.Context, id string) error
}
