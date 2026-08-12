package main

import (
	"fmt"
	"log"

	"github.com/fastsprint/common/config"
	"github.com/fastsprint/common/database"
	"github.com/fastsprint/common/middleware"
	"github.com/fastsprint/project-service/internal/handler"
	"github.com/fastsprint/project-service/internal/repository"
	"github.com/fastsprint/project-service/internal/service"
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

	projectRepo := repository.NewProjectRepository(database.DB)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		projects := api.Group("/projects")
		{
			projects.POST("", projectHandler.CreateProject)
			projects.GET("", projectHandler.ListProjects)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)
			projects.GET("/key/:key", projectHandler.GetProjectByKey)

			members := projects.Group("/:id/members")
			{
				members.GET("", projectHandler.ListMembers)
				members.POST("", projectHandler.AddMember)
				members.DELETE("/:userId", projectHandler.RemoveMember)
				members.PUT("/:userId/role", projectHandler.UpdateMemberRole)
			}

			statuses := projects.Group("/:id/statuses")
			{
				statuses.GET("", projectHandler.ListStatuses)
				statuses.POST("", projectHandler.CreateStatus)
			}

			issueTypes := projects.Group("/:id/issue-types")
			{
				issueTypes.GET("", projectHandler.ListIssueTypes)
				issueTypes.POST("", projectHandler.CreateIssueType)
			}

			priorities := projects.Group("/:id/priorities")
			{
				priorities.GET("", projectHandler.ListPriorities)
				priorities.POST("", projectHandler.CreatePriority)
			}
		}
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Project service starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
