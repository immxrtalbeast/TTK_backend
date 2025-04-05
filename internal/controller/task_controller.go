package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/TTK_backend/internal/domain"
)

type TaskController struct {
	interactor domain.TaskInteractor
}

func NewTaskController(interactor domain.TaskInteractor) *TaskController {
	return &TaskController{interactor: interactor}
}

func (c *TaskController) Task(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing task ID"})
		return
	}
	Task, err := c.interactor.Task(ctx, idStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get task",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Task": Task,
	})

}

func (c *TaskController) CreateTask(ctx *gin.Context) {
	type CreateTaskRequest struct {
		Title     string          `json:"title" binding:"required,min=3,max=50"`
		Image     string          `json:"image" binding:"required"`
		Content   string          `json:"content"`
		UserID    string          `json:"user_id"`
		PlannedAt time.Time       `json:"planned_at"`
		Priority  domain.Priority `json:"priority" binding:"required"`
		Status    domain.Status   `json:"status"`
	}
	var req CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	taskID, err := c.interactor.CreateTask(ctx, req.Title, req.Image, req.Content, req.UserID, req.Priority, req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to update create task",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"taskID": taskID,
		"userID": ctx.Keys["userName"],
	})

}
