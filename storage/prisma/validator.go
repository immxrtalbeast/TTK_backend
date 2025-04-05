package prisma

import (
	"github.com/immxrtalbeast/TTK_backend/internal/domain"
	"github.com/immxrtalbeast/TTK_backend/storage/prisma/db"
)

func ValidateUser(userDB db.UserModel) domain.User {
	user := domain.User{
		ID:       userDB.ID,
		Name:     userDB.FullName,
		Login:    userDB.Login,
		PassHash: []byte(userDB.PasswordHash),
		IsAdmin:  domain.Role(userDB.Role),
	}
	return user
}

func ValidateArticle(articleDB db.ArticleModel) domain.Article {
	updated, _ := articleDB.UpdatedAt()
	nameLastReedit, _ := articleDB.NameLastReedit()
	content, _ := articleDB.Content()
	atricle := domain.Article{
		ID:               articleDB.ID,
		Title:            articleDB.Title,
		UpdatedAt:        updated,
		CreatedAt:        articleDB.CreatedAt,
		Name_last_reedit: nameLastReedit,
		Image:            articleDB.Image,
		Content:          content,
	}
	return atricle

}

func ValidateTask(taskDB db.TaskModel) domain.Task {
	task := domain.Task{
		ID:        taskDB.ID,
		Title:     taskDB.Title,
		UserID:    taskDB.UserID,
		PlannedAt: taskDB.PlannedAt,
		CreatedAt: taskDB.CreatedAt,
		Priority:  domain.Priority(taskDB.Priority),
		Status:    domain.Status(taskDB.Status),
	}
	return task
}
