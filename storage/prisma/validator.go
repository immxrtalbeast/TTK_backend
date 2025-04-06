package prisma

import (
	"github.com/immxrtalbeast/TTK_backend/internal/domain"
	"github.com/immxrtalbeast/TTK_backend/storage/prisma/db"
)

func ValidateUser(userDB db.UserModel) domain.User {
	user := domain.User{
		ID:        userDB.ID,
		Name:      userDB.FullName,
		Login:     userDB.Login,
		PassHash:  []byte(userDB.PasswordHash),
		CreatedAt: userDB.CreatedAt,
		IsAdmin:   domain.Role(userDB.Role),
	}
	return user
}

func ValidateArticle(articleDB db.ArticleModel) domain.Article {
	content, _ := articleDB.Content()
	atricle := domain.Article{
		ID:         articleDB.ID,
		Title:      articleDB.Title,
		UpdatedAt:  articleDB.UpdatedAt,
		CreatedAt:  articleDB.CreatedAt,
		Creator:    articleDB.CreatorName,
		LastEditor: articleDB.LastEditorName,
		Image:      articleDB.Image,
		Content:    content,
	}
	return atricle

}

//	ID           string
// Title        string
// ArticleId    string
// UserId       string
// ChangedAt    time.Time
// EventType    EventType
// ArticleTitle string

func ValidateArticleHistory(historyDB db.ArticleHistoryModel) domain.History {
	history := domain.History{
		ID:           historyDB.ID,
		UserId:       historyDB.UserID,
		ArticleId:    historyDB.ArticleID,
		ChangedAt:    historyDB.ChangedAt,
		EventType:    domain.EventType(historyDB.EventType),
		ArticleTitle: historyDB.ArticleTitle,
	}
	return history
}

func ValidateTask(taskDB db.TaskModel) domain.Task {
	respUser := taskDB.Responsibleuser()
	user := ValidateUser(*respUser)
	task := domain.Task{
		ID:               taskDB.ID,
		Title:            taskDB.Title,
		ReliableUserName: user.Name,
		UserID:           taskDB.UserID,
		PlannedAt:        taskDB.PlannedAt,
		CreatedAt:        taskDB.CreatedAt,
		Priority:         domain.Priority(taskDB.Priority),
		Status:           domain.Status(taskDB.Status),
	}
	return task
}
