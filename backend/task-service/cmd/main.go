package main

import (
	"fmt"
	"log"

	"github.com/fastsprint/common/config"
	"github.com/fastsprint/common/database"
	"github.com/fastsprint/common/middleware"
	"github.com/fastsprint/task-service/internal/handler"
	"github.com/fastsprint/task-service/internal/repository"
	"github.com/fastsprint/task-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)

	_, err = database.InitPostgres(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	_, err = database.InitRedis(cfg.Redis)
	if err != nil {
		log.Printf("Warning: Failed to init redis: %v", err)
	}

	taskRepo := repository.NewTaskRepository(database.DB)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		issues := api.Group("/issues")
		{
			issues.POST("", taskHandler.CreateIssue)
			issues.GET("", taskHandler.ListIssues)
			issues.GET("/:id", taskHandler.GetIssue)
			issues.PUT("/:id", taskHandler.UpdateIssue)
			issues.DELETE("/:id", taskHandler.DeleteIssue)
			issues.GET("/key/:key", taskHandler.GetIssueByKey)

			comments := issues.Group("/:id/comments")
			{
				comments.GET("", taskHandler.GetComments)
				comments.POST("", taskHandler.AddComment)
				comments.DELETE("/:commentId", taskHandler.DeleteComment)
			}

			changelogs := issues.Group("/:id/changelogs")
			{
				changelogs.GET("", taskHandler.GetChangelogs)
			}
		}

		sprints := api.Group("/sprints")
		{
			sprints.POST("", taskHandler.CreateSprint)
			sprints.GET("", taskHandler.ListSprints)
			sprints.GET("/:id", taskHandler.GetSprint)
			sprints.PUT("/:id", taskHandler.UpdateSprint)
			sprints.DELETE("/:id", taskHandler.DeleteSprint)
			sprints.POST("/:id/start", taskHandler.StartSprint)
			sprints.POST("/:id/complete", taskHandler.CompleteSprint)
		}

		board := api.Group("/board")
		{
			board.GET("", taskHandler.GetBoardData)
		}
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Task service starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
