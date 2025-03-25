package routes

import (
	"github.com/Andrew44Ashraf/fintech-service/internal/handlers"
	"github.com/Andrew44Ashraf/fintech-service/internal/repository"
	"github.com/gin-gonic/gin"
	"database/sql"
)

// SetupRoutes initializes all API routes with dependency injection
func SetupRoutes(router *gin.Engine, db *sql.DB) {
	// Initialize repositories
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize handlers
	accountHandler := handlers.NewAccountHandler(accountRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionRepo, accountRepo)

	// API routes
	api := router.Group("/api")
	{
		// Account routes
		api.POST("/accounts", accountHandler.OpenAccount)
		api.GET("/accounts/:id/balance", accountHandler.GetBalance)

		// Transaction routes
		api.POST("/accounts/:id/deposit", transactionHandler.Deposit)
		api.POST("/accounts/:id/withdraw", transactionHandler.Withdraw)
		api.GET("/accounts/:id/transactions", transactionHandler.GetTransactions) // ?limit=10&offset=0
	}
}