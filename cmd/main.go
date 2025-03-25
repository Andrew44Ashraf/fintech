package main

import (
	"database/sql"
	"github.com/Andrew44Ashraf/fintech-service/internal/database"
	"github.com/Andrew44Ashraf/fintech-service/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize DB
	db := database.Connect()
	defer db.Close()

	// Create router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, db)

	// Start server
	router.Run(":8080")
}