package main

import (
	"log/slog"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/TTK_backend/internal/config"
	"github.com/immxrtalbeast/TTK_backend/internal/controller"
	"github.com/immxrtalbeast/TTK_backend/internal/middleware"
	"github.com/immxrtalbeast/TTK_backend/internal/usecase/article"
	"github.com/immxrtalbeast/TTK_backend/internal/usecase/task"
	"github.com/immxrtalbeast/TTK_backend/internal/usecase/user"
	"github.com/immxrtalbeast/TTK_backend/storage/prisma"
	"github.com/joho/godotenv"
)

// go run cmd/main.go --config=./config/local.yaml

// perpage articles, tasks,
func main() {
	cfg := config.MustLoad()
	log := setupLogger()
	log.Info("starting application", slog.Any("config", cfg))
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	db, err := prisma.New()
	if err != nil {
		panic("Failed to connect DB" + err.Error())
	}
	defer db.Disconnect()
	//TODO: validate data, implement more methods. think about history
	userINT := user.NewUserInteractor(db, cfg.TokenTTL, cfg.AppSecret)
	userController := controller.NewUserController(userINT)

	articleINT := article.NewArticleInteractor(db)
	articleController := controller.NewArticleController(articleINT, log)

	taskINT := task.NewTaskInteractor(db)
	taskController := controller.NewTaskController(taskINT)

	authMiddleware := middleware.AuthMiddleware(cfg.AppSecret)
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	router.Use(cors.New(config))
	api := router.Group("/api/v1")
	{
		article := api.Group("/article")
		article.Use(authMiddleware)
		{
			article.POST("/create", articleController.CreateArticle)
			article.GET("/:id", articleController.Article)
			article.GET("/get_page", articleController.Articles)
			article.DELETE("/:id", articleController.DeleteArticle)
		}
		task := api.Group("/task")
		// task.Use(authMiddleware)
		{
			task.POST("/create", taskController.CreateTask)
			task.GET("/:id", taskController.Task)
		}
		api.POST("/register", userController.CreateUser)
		api.GET("/user/:id", userController.User)
		api.POST("/login", userController.Login)
	}
	router.Run(":8080")

}
func setupLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	return log
}
