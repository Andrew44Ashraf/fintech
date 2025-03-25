package main

import (
    "github.com/Andrew44Ashraf/fintech-service/internal/database"
    "github.com/Andrew44Ashraf/fintech-service/internal/routes"
    "github.com/gin-gonic/gin"
    "github.com/swaggo/files"
    "github.com/swaggo/gin-swagger"
    "log"
)

// @title Fintech API
// @version 1.0
// @description API for managing accounts and transactions in a fintech application.
// @host localhost:8080
// @BasePath /api

func main() {
	database.InitDB()

	r := gin.Default()

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	routes.SetupRoutes(r)

	log.Println("Server is running on port 8080...")
	r.Run(":8080")
}
