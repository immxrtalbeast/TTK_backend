package controller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/TTK_backend/internal/domain"
)

type ArticleController struct {
	interactor domain.ArticleInteractor
	log        *slog.Logger
}

func NewArticleController(interactor domain.ArticleInteractor, log *slog.Logger) *ArticleController {
	return &ArticleController{interactor: interactor, log: log}
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

func (c *ArticleController) CreateArticle(ctx *gin.Context) {
	type CreateArticleRequest struct {
		Title   string `json:"title" binding:"required,min=3,max=50"`
		Image   string `json:"image" binding:"required"`
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
	userName, _ := ctx.Keys["userName"].(string)
	article, err := c.interactor.CreateArticle(ctx, req.Title, req.Image, req.Content, userName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create article",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"article": article,
	})

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
	article, err := c.interactor.UpdateArticle(ctx, req.Title, req.Image, req.Content, userName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to update article",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"article": article,
	})
}
