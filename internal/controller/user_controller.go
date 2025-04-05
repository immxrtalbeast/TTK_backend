package controller

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/TTK_backend/internal/domain"
)

type UserController struct {
	interactor domain.UserInteractor
}

func NewUserController(interactor domain.UserInteractor) *UserController {
	return &UserController{interactor: interactor}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	type CreateUserRequest struct {
		Login string `json:"login" binding:"required,min=3,max=50"`
		Name  string `json:"name" binding:"required"`
		Pass  string `json:"password" binding:"required"`
	}

	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Валидация логина
	loginRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	if !loginRegex.MatchString(req.Login) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid login",
			"details": "Login must contain only latin letters",
		})
		return
	}

	// Валидация имени
	nameRegex := regexp.MustCompile(`^[а-яА-ЯёЁ\s]+$`)
	if !nameRegex.MatchString(req.Name) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid name",
			"details": "Name must contain only russian letters and spaces",
		})
		return
	}

	// Валидация пароля
	passRegex := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+\[\]{};:<>,./?~\\-]+$`)
	if !passRegex.MatchString(req.Pass) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid password",
			"details": "Password contains forbidden characters",
		})
		return
	}

	// Если все проверки пройдены
	if err := c.interactor.CreateUser(ctx, req.Login, req.Name, req.Pass); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create user",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *UserController) User(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing article ID"})
		return
	}
	user, err := c.interactor.User(ctx, idStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get user",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (c *UserController) Login(ctx *gin.Context) {
	type LoginRequest struct {
		Login string `json:"login" binding:"required"`
		Pass  string `json:"password" binding:"required"`
	}
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}
	token, err := c.interactor.Login(ctx, req.Login, req.Pass)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to login",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}
