package controller

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"github.com/immxrtalbeast/TTK_backend/internal/domain"
)

type UserController struct {
	interactor domain.UserInteractor
}

func NewUserController(interactor domain.UserInteractor) *UserController {
	return &UserController{interactor: interactor}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	if validator, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validator.RegisterValidation("russian_alpha", validateRussianAlpha)
		validator.RegisterValidation("password", validatePassword)
	}

	type CreateUserRequest struct {
		Login string `json:"login" binding:"required,min=3,max=50,alpha"`
		Name  string `json:"name" binding:"required,russian_alpha"`
		Pass  string `json:"password" binding:"required,min=8,password"`
	}
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleValidationErrors(ctx, err)
		return
	}
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

// Валидатор для русского алфавита и пробелов
func validateRussianAlpha(fl validator.FieldLevel) bool {
	matched, _ := regexp.MatchString(`^[а-яА-ЯёЁ\s-]+$`, fl.Field().String())
	return matched
}

// Валидатор для пароля (латинские буквы, цифры, специальные символы)
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Регулярное выражение: минимум 1 буква, 1 цифра, 1 спецсимвол
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString
	hasLower := regexp.MustCompile(`[a-z]`).MatchString
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+~-]`).MatchString

	return (hasUpper(password) || hasLower(password)) &&
		hasNumber(password) &&
		hasSpecial(password) &&
		regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+~-]+$`).MatchString(password)
}

func handleValidationErrors(ctx *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errors := make(map[string]string)
		for _, e := range ve {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = "This field is required"
			case "min":
				errors[field] = fmt.Sprintf("Minimum length is %s", e.Param())
			case "max":
				errors[field] = fmt.Sprintf("Maximum length is %s", e.Param())
			case "alpha":
				errors[field] = "Should contain only latin letters"
			case "russian_alpha":
				errors[field] = "Should contain only russian letters"
			case "password":
				errors[field] = "Must include letters, numbers and special symbols"
			default:
				errors[field] = "Invalid value"
			}
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
