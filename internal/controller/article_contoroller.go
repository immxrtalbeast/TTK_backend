package controller

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/TTK_backend/internal/domain"
)

type ArticleController struct {
	interactor  domain.ArticleInteractor
	hInteractor domain.HistoryInteractor
	log         *slog.Logger
}

func NewArticleController(interactor domain.ArticleInteractor, hinteractor domain.HistoryInteractor, log *slog.Logger) *ArticleController {
	return &ArticleController{interactor: interactor, hInteractor: hinteractor, log: log}
}

func (c *ArticleController) Article(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing article ID"})
		return
	}
	article, err := c.interactor.Article(ctx, idStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get article",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"article": article,
	})

}

func (c *ArticleController) Articles(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("p", "1")
	limitStr := ctx.DefaultQuery("limit", "6")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	articles, err := c.interactor.Articles(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{
			"error":   "failed to create article",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"articles": articles,
	})
}

func (c *ArticleController) CreateArticle(ctx *gin.Context) {
	type CreateArticleRequest struct {
		Title   string `json:"title" binding:"required,min=3,max=50"`
		Image   string `json:"image"`
		Content string `json:"content"`
	}
	var req CreateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}
	userID, _ := ctx.Keys["userID"].(string)
	article, err := c.interactor.CreateArticle(ctx, req.Title, req.Image, req.Content, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create article",
			"details": err.Error(),
		})
		return
	}
	err = c.hInteractor.InitHistory(ctx, article.ID, userID, req.Title)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create history",
			"details": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"article": article,
	})

}
func (c *ArticleController) DeleteArticle(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing article ID"})
		return
	}
	userID, _ := ctx.Keys["userID"].(string)
	article, _ := c.interactor.Article(ctx, id)
	err := c.interactor.DeteleArticle(ctx, id, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to delete article",
			"details": err.Error(),
		})
		return
	}
	c.hInteractor.UpdateHistory(ctx, id, userID, domain.EventType("DELETE"), article.Title)
	ctx.JSON(http.StatusOK, gin.H{})
}
func (c *ArticleController) UpdateArticle(ctx *gin.Context) {
	type UpdateArticleRequest struct {
		ID      string `json:"id" binding:"required"`
		Title   string `json:"title" binding:"required,min=3,max=50"`
		Image   string `json:"image" binding:"required"`
		Content string `json:"content"`
	}
	var req UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}
	userName, _ := ctx.Keys["userName"].(string)
	userID, _ := ctx.Keys["userID"].(string)
	article, err := c.interactor.UpdateArticle(ctx, req.ID, req.Title, req.Image, req.Content, userName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to update article",
			"details": err.Error(),
		})
		return
	}
	c.hInteractor.UpdateHistory(ctx, article.ID, userID, domain.EventType("UPDATED"), article.Title)
	ctx.JSON(http.StatusOK, gin.H{
		"article": article,
	})
}
