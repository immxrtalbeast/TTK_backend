package prisma

import (
	"context"
	"fmt"

	"github.com/immxrtalbeast/TTK_backend/internal/domain"
	"github.com/immxrtalbeast/TTK_backend/storage/prisma/db"
)

type Storage struct {
	client *db.PrismaClient
}

func New() (*Storage, error) {
	const op = "storage.prisma.New"

	client := db.NewClient()

	if err := client.Prisma.Connect(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{client: client}, nil
}
func (s *Storage) Disconnect() {
	s.client.Prisma.Disconnect()
}

// USER
func (s *Storage) User(ctx context.Context, id string) (*domain.User, error) {
	const op = "storage.user.get"
	userDB, err := s.client.User.FindUnique(db.User.ID.Equals(id)).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	user := ValidateUser(*userDB)
	return &user, nil
}
func (s *Storage) UserByLogin(ctx context.Context, login string) (*domain.User, error) {
	const op = "storage.user.get_by_login"
	userDB, err := s.client.User.FindUnique(db.User.Login.Equals(login)).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	user := ValidateUser(*userDB)
	return &user, nil
}
func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
	const op = "storage.user.create"
	_, err := s.client.User.CreateOne(
		db.User.Login.Set(user.Login),
		db.User.PasswordHash.Set(string(user.PassHash)),
		db.User.FullName.Set(user.Name),
	).Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// ARTICLE
func (s *Storage) CreateArticle(ctx context.Context, article *domain.Article) (string, error) {
	const op = "storage.create.article"
	result, err := s.client.Article.CreateOne(
		db.Article.Title.Set(article.Title),
		db.Article.Image.Set(article.Image),
		db.Article.UpdatedAt.Set(article.UpdatedAt),
		db.Article.NameLastReedit.Set(article.Name_last_reedit),
		db.Article.Content.Set(article.Content),
	).Exec(ctx)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return result.ID, nil
}

func (s *Storage) Article(ctx context.Context, id string) (*domain.Article, error) {
	const op = "storage.article.get"
	articleDB, err := s.client.Article.FindUnique(
		db.Article.ID.Equals(id),
	).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("article not found")
	}
	article := ValidateArticle(*articleDB)
	return &article, nil
}

//TASK

func (s *Storage) CreateTask(ctx context.Context, task *domain.Task) (string, error) {
	const op = "storage.task.create"
	user := db.User.ID.Equals(task.UserID)
	result, err := s.client.Task.CreateOne(
		db.Task.Title.Set(task.Title),
		db.Task.User.Link(user),
		db.Task.PlannedAt.Set(task.PlannedAt),
		db.Task.Priority.Set(db.Priority(task.Priority)),
		db.Task.Status.Set(db.Status(task.Status)),
	).Exec(ctx)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return result.ID, nil
}

func (s *Storage) Task(ctx context.Context, id string) (*domain.Task, error) {
	const op = "storage.task.get"
	taskDB, err := s.client.Task.FindUnique(db.Task.ID.Equals(id)).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}
	task := ValidateTask(*taskDB)
	return &task, nil

}
