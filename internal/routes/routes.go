package routes

import (
	"github.com/Andrew44Ashraf/fintech-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes defines API routes
func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/accounts", handlers.OpenAccount)
		api.GET("/accounts/:id/balance", handlers.GetBalance)
		api.POST("/accounts/:id/deposit", handlers.Deposit)
		api.POST("/accounts/:id/withdraw", handlers.Withdraw)
		api.GET("/accounts/:id/transactions", handlers.GetTransactions) // Supports ?limit=10&offset=0
	}
}
