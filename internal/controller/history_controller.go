package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/TTK_backend/internal/domain"
)

type HistoryController struct {
	interactor domain.HistoryInteractor
}

func NewHistoryController(interactor domain.HistoryInteractor) *HistoryController {
	return &HistoryController{interactor: interactor}
}

func (c *HistoryController) History(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("p", "1")
	limitStr := ctx.DefaultQuery("limit", "6")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	histories, err := c.interactor.Histories(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get history",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": histories,
	})
}
