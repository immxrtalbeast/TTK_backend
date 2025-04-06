package prisma

import (
	"context"
	"fmt"
	"time"

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

func (s *Storage) Users(ctx context.Context, page int, limit int) ([]*domain.User, error) {
	const op = "storage.users.all"
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 6 // значение по умолчанию
	}
	skip := (page - 1) * limit

	// Добавляем пагинацию и сортировку
	usersDB, err := s.client.User.FindMany().
		Take(limit). // Количество элементов на странице
		Skip(skip).
		OrderBy(db.User.CreatedAt.Order(db.ASC)). // Сколько элементов пропуститьA
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var users []*domain.User
	for _, userDB := range usersDB {
		user := ValidateUser(userDB)
		users = append(users, &user)
	}
	return users, nil

}

func (s *Storage) UpdateUser(ctx context.Context, user *domain.User) error {
	const op = "storage.user.update"
	_, err := s.client.User.FindUnique(db.User.ID.Equals(user.ID)).Update(
		db.User.Login.Set(user.Login),
		db.User.PasswordHash.Set(string(user.PassHash)),
		db.User.FullName.Set(user.Name),
		db.User.Role.Set(db.Role(user.IsAdmin)),
	).Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
func (s *Storage) DeleteUser(ctx context.Context, id string) error {
	const op = "storage.user.delete"
	_, err := s.client.User.FindUnique(db.User.ID.Equals(id)).Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// ARTICLE
func (s *Storage) CreateArticle(ctx context.Context, article *domain.Article) (string, error) {
	const op = "storage.article.create"
	result, err := s.client.Article.CreateOne(
		db.Article.Title.Set(article.Title),
		db.Article.LastEditorName.Set(article.LastEditor),
		db.Article.CreatorName.Set(article.LastEditor),
		db.Article.Image.Set(article.Image),
		db.Article.UpdatedAt.Set(article.UpdatedAt),
		db.Article.Content.Set(article.Content),
	).Exec(ctx)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return result.ID, nil
}
func (s *Storage) Articles(ctx context.Context, page, limit int) ([]*domain.Article, error) {
	const op = "storage.article.get_all"
	// Валидация параметров пагинации
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 6 // значение по умолчанию
	}
	skip := (page - 1) * limit

	// Добавляем пагинацию и сортировку
	articlesDB, err := s.client.Article.FindMany().
		Take(limit). // Количество элементов на странице
		Skip(skip).
		OrderBy(db.Article.CreatedAt.Order(db.ASC)). // Сколько элементов пропуститьA
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var articles []*domain.Article
	for _, articleDB := range articlesDB {
		atricle := ValidateArticle(articleDB)
		articles = append(articles, &atricle)
	}
	return articles, nil

}

func (s *Storage) Article(ctx context.Context, id string) (*domain.Article, error) {
	articleDB, err := s.client.Article.FindUnique(
		db.Article.ID.Equals(id),
	).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("article not found")
	}
	article := ValidateArticle(*articleDB)
	return &article, nil
}

func (s *Storage) UpdateArticle(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	const op = "storage.article.update"
	articleDB, err := s.client.Article.FindUnique(db.Article.ID.Equals(article.ID)).Update(
		db.Article.Title.Set(article.Title),
		db.Article.UpdatedAt.Set(time.Now()),
		db.Article.LastEditorName.Set(article.LastEditor),
		db.Article.Image.Set(article.Image),
		db.Article.Content.Set(article.Content),
	).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	articleDomain := ValidateArticle(*articleDB)
	return &articleDomain, nil
}

func (s *Storage) DeleteArticle(ctx context.Context, id string) error {
	const op = "storage.article.delete"
	_, err := s.client.Article.FindUnique(db.Article.ID.Equals(id)).Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// HISTORY

func (s *Storage) InitHistory(ctx context.Context, history *domain.History) error {
	const op = "storage.history.init"
	_, err := s.client.ArticleHistory.CreateOne(
		db.ArticleHistory.ArticleID.Set(history.ArticleId),
		db.ArticleHistory.UserID.Set(history.UserId),
		db.ArticleHistory.EventType.Set(db.EventTypeCreate),
		db.ArticleHistory.ArticleTitle.Set(history.ArticleTitle),
	).Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) UpdateHistory(ctx context.Context, history *domain.History) error {
	const op = "storage.history.init"
	_, err := s.client.ArticleHistory.CreateOne(
		db.ArticleHistory.ArticleID.Set(history.ArticleId),
		db.ArticleHistory.UserID.Set(history.UserId),
		db.ArticleHistory.EventType.Set(db.EventType(history.EventType)),
		db.ArticleHistory.ArticleTitle.Set(history.ArticleTitle),
	).Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) Histories(ctx context.Context, page, limit int) ([]*domain.History, error) {
	const op = "storage.history.all"

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 6 // значение по умолчанию
	}
	skip := (page - 1) * limit
	historiesDB, err := s.client.ArticleHistory.FindMany().
		Take(limit). // Количество элементов на странице
		Skip(skip).
		OrderBy(db.ArticleHistory.ChangedAt.Order(db.ASC)). // Сколько элементов пропуститьA
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var histories []*domain.History
	for _, historyDB := range historiesDB {
		history := ValidateArticleHistory(historyDB)
		histories = append(histories, &history)
	}
	return histories, nil
}

//TASK

func (s *Storage) CreateTask(ctx context.Context, task *domain.Task) (string, error) {
	const op = "storage.task.create"
	user := db.User.ID.Equals(task.UserID)
	result, err := s.client.Task.CreateOne(
		db.Task.Title.Set(task.Title),
		db.Task.Content.Set(task.Content),
		db.Task.Image.Set(task.Image),
		db.Task.Responsibleuser.Link(user),
		db.Task.PlannedAt.Set(task.PlannedAt),
		db.Task.Priority.Set(db.Priority(task.Priority)),
		db.Task.Status.Set(db.Status(task.Status)),
	).Exec(ctx)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return result.ID, nil
}

func (s *Storage) Tasks(ctx context.Context, page, limit int) ([]*domain.Task, error) {
	const op = "storage.task.all"

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 6 // значение по умолчанию
	}
	skip := (page - 1) * limit
	tasksDB, err := s.client.Task.FindMany().
		Take(limit). // Количество элементов на странице
		Skip(skip).
		OrderBy(db.Task.CreatedAt.Order(db.ASC)). // Сколько элементов пропуститьA
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var tasks []*domain.Task
	for _, taskDB := range tasksDB {
		task := ValidateTask(taskDB)
		tasks = append(tasks, &task)
	}
	return tasks, nil
}
func (s *Storage) Task(ctx context.Context, id string) (*domain.Task, error) {
	taskDB, err := s.client.Task.FindUnique(db.Task.ID.Equals(id)).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}
	task := ValidateTask(*taskDB)
	return &task, nil

}
func (s *Storage) UpdateTask(ctx context.Context, task *domain.Task) error {
	const op = "storage.task.update"
	_, err := s.client.Task.FindUnique(db.Task.ID.Equals(task.ID)).Update(
		db.Task.Title.Set(task.Title),
		db.Task.Content.Set(task.Content),
		db.Task.Image.Set(task.Image),
		db.Task.UserID.Set(task.UserID),
		db.Task.PlannedAt.Set(task.PlannedAt),
		db.Task.Priority.Set(db.Priority(task.Priority)),
		db.Task.Status.Set(db.Status(task.Status)),
	).Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

}

func (s *Storage) DeleteTask(ctx context.Context, id string) error {
	const op = "storage.task.delete"
	_, err := s.client.Task.FindUnique(db.Task.ID.Equals(id)).Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
