package main

import (
	"expensetracker/config"
	"expensetracker/docs"
	handlers "expensetracker/handlers"
	"expensetracker/middleware"
	"expensetracker/repository"
	"expensetracker/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Expense Tracker API
// @version 1.0
// @description RESTful API for tracking expenses
// @termsOfService http://example.com/terms/

// @contact.name Bhavani
// @contact.email bhavani@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8085
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg := config.LoadConfig()

	db := config.ConnectDatabase(cfg)

	router := gin.Default()
	router.Use(gin.Logger(), gin.Recovery(), middleware.RateLimiter())

	expenseRepo := repository.NewExpenseRepository(db)
	expenseService := services.NewExpenseService(expenseRepo)
	expenseHandler := handlers.NewExpenseHandler(expenseService)

	authService := services.NewAuthService(cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	api := router.Group("/api/v1")
	{
		api.POST("/auth/login", authHandler.Login)

		expense := api.Group("/expenses")
		expense.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			expense.POST("", expenseHandler.CreateExpense)
			expense.GET("/:id", expenseHandler.GetExpense)
			expense.PUT("/:id", expenseHandler.UpdateExpense)
			expense.DELETE("/:id", expenseHandler.DeleteExpense)
			expense.GET("", expenseHandler.ListExpenses)
			expense.GET("/summary", expenseHandler.SummaryExpenses)
		}
	}

	docs.SwaggerInfo.BasePath = "/"
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8085/swagger/doc.json"),
	))

	log.Println("Starting server on port 8085...")
	if err := router.Run(":8085"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
