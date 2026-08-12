package main

import (
	"fmt"
	"log"

	"github.com/fastsprint/common/config"
	"github.com/fastsprint/common/database"
	"github.com/fastsprint/common/middleware"
	"github.com/fastsprint/user-service/internal/handler"
	"github.com/fastsprint/user-service/internal/repository"
	"github.com/fastsprint/user-service/internal/service"
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

	userRepo := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		user := api.Group("/users")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/me", userHandler.GetCurrentUser)
			user.PUT("/me", userHandler.UpdateCurrentUser)
			user.POST("/change-password", userHandler.ChangePassword)
			user.GET("/:id", userHandler.GetUserByID)
			user.GET("", userHandler.ListUsers)
		}
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("User service starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
