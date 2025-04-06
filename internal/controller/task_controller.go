package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/TTK_backend/internal/domain"
)

type TaskController struct {
	interactor  domain.TaskInteractor
	hInteractor domain.HistoryInteractor
}

func NewTaskController(interactor domain.TaskInteractor, hInteractor domain.HistoryInteractor) *TaskController {
	return &TaskController{interactor: interactor, hInteractor: hInteractor}
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

func (c *TaskController) Tasks(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("p", "1")
	limitStr := ctx.DefaultQuery("limit", "6")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	tasks, err := c.interactor.Tasks(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get tasks",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})

}

func (c *TaskController) UpdateTask(ctx *gin.Context) {
	type UpdateTaskRequest struct {
		ID        string          `json:"id" binding:"required"`
		Title     string          `json:"title" binding:"required,min=3,max=50"`
		Image     string          `json:"image" binding:"required"`
		Content   string          `json:"content"`
		UserID    string          `json:"user_id"`
		PlannedAt time.Time       `json:"planned_at"`
		Priority  domain.Priority `json:"priority" binding:"required"`
		Status    domain.Status   `json:"status"`
	}
	var req UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}
	err := c.interactor.UpdateTask(ctx, req.ID, req.Title, req.Image, req.Content, req.UserID, req.Priority, req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to update task",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})

}

func (c *TaskController) CreateTask(ctx *gin.Context) {
	type CreateTaskRequest struct {
		Title     string          `json:"title" binding:"required,min=3,max=50"`
		Image     string          `json:"image"`
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

	taskID, err := c.interactor.CreateTask(ctx, req.Title, req.Image, req.Content, req.PlannedAt, req.UserID, req.Priority, req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create task",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"taskID": taskID,
		"userID": ctx.Keys["userName"],
	})

}

func (c *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing article ID"})
		return
	}

	err := c.interactor.DeleteTask(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create task",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
