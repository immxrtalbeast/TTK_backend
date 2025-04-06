package task

import (
	"context"
	"errors"
	"fmt"

	"github.com/immxrtalbeast/TTK_backend/internal/domain"
	"github.com/immxrtalbeast/TTK_backend/storage/prisma/db"
)

type TaskInteractor struct {
	taskRepo domain.TaskRepository
}

func NewTaskInteractor(taskRepo domain.TaskRepository) domain.TaskInteractor {
	return &TaskInteractor{taskRepo: taskRepo}
}

func (ai *TaskInteractor) Task(ctx context.Context, id string) (*domain.Task, error) {
	const op = "uc.task.get"

	task, err := ai.taskRepo.Task(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, db.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)

	}
	return task, nil
}

func (ai *TaskInteractor) Tasks(ctx context.Context, page, limit int) ([]*domain.Task, error) {
	const op = "uc.tast.all"
	tasks, err := ai.taskRepo.Tasks(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tasks, nil
}

func (ai *TaskInteractor) UpdateTask(ctx context.Context, id string, title string, image string, content string, userID string, priority domain.Priority, status domain.Status) error {
	const op = "uc.task.update"
	task := domain.Task{
		ID:       id,
		Title:    title,
		Image:    image,
		Content:  content,
		UserID:   userID,
		Priority: priority,
		Status:   status,
	}
	err := ai.taskRepo.UpdateTask(ctx, &task)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (ai *TaskInteractor) DeleteTask(ctx context.Context, id string) error {
	const op = "uc.task.delete"
	err := ai.taskRepo.DeleteTask(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (ai *TaskInteractor) CreateTask(ctx context.Context, title string, image string, content string, userID string, priority domain.Priority, status domain.Status) (string, error) {
	const op = "uc.task.create"
	task := domain.Task{
		Title:    title,
		Image:    image,
		Content:  content,
		UserID:   userID,
		Priority: priority,
		Status:   status,
	}

	taskID, err := ai.taskRepo.CreateTask(ctx, &task)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return taskID, nil

}
