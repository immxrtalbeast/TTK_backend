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

// func (s *Storage) InitHistory(ctx context.Context, article *domain.Article) error {
// 	const op = "storage.history.init"
// 	user := db.User.ID.Equals(article.Name_last_reedit)
// 	_, err := s.client.ArticleHistory.CreateOne(
// 		db.ArticleHistory.ArticleID.Set(article.ID),
// 		db.ArticleHistory.User.Link(user),
// 		db.ArticleHistory.EventType.Set("CREATE"),
// 		db.ArticleHistory.ArticleTitle.Set(article.Title),
// 		db.ArticleHistory.ChangedAt.Set(time.Now()),
// 	).Exec(ctx)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}
// 	return nil
// }

// func (s *Storage) HistoryByID(ctx context.Context, id string) (*domain.History, error) {
// 	const op = "storage.history.by_id"
// 	historyDB, err := s.client.ArticleHistory.FindUnique(db.ArticleHistory.ID.Equals(id)).Exec(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	history := ValidateArticleHistory(*historyDB)
// 	return &history, nil
// }

// func (s *Storage) Histories(ctx context.Context) ([]*domain.History, error) {
// 	const op = "storage.history.all"
// 	historiesDB, err := s.client.ArticleHistory.FindMany().Exec(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	var histories []*domain.History
// 	for _, historyDB := range historiesDB {
// 		history := ValidateArticleHistory(historyDB)
// 		histories = append(histories, &history)
// 	}
// 	return histories, nil
// }

// func (s *Storage) HistoryByArticleID(ctx context.Context, article_id string) ([]*domain.History, error) {
// 	const op = "storage.history.by_article_id"
// 	historiesDB, err := s.client.ArticleHistory.FindMany(db.ArticleHistory.ArticleID.Equals(article_id)).Exec(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	var histories []*domain.History
// 	for _, historyDB := range historiesDB {
// 		history := ValidateArticleHistory(historyDB)
// 		histories = append(histories, &history)
// 	}
// 	return histories, nil
// }

// func (s *Storage) AddHistory(ctx context.Context, articleID string, userID string, eventType domain.EventType, articleTitle string) error {
// 	const op = "storage.history.update"
// 	user := db.User.ID.Equals(userID)
// 	_, err := s.client.ArticleHistory.CreateOne(
// 		db.ArticleHistory.ArticleID.Set(articleID),
// 		db.ArticleHistory.User.Link(user),
// 		db.ArticleHistory.EventType.Set(db.EventType(eventType)),
// 		db.ArticleHistory.ArticleTitle.Set(articleTitle),
// 	).Exec(ctx)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}
// 	return nil
// }

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
