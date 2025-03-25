package handlers

import (
	"fintech-service/internal/database"
	"fintech-service/internal/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// OpenAccount godoc
// @Summary Create a new account
// @Description Opens a new account with a default balance of 0
// @Tags accounts
// @Accept json
// @Produce json
// @Success 200 {object} map[string]int "Returns account ID"
// @Failure 500 {object} map[string]string "Server error"
// @Router /accounts [post]
func OpenAccount(c *gin.Context) {
	accountID, err := repository.CreateAccount(database.DB)
	if err != nil {
		log.Println("Failed to create account:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"account_id": accountID})
}

// GetBalance godoc
// @Summary Get account balance
// @Description Returns the current balance of an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} map[string]float64 "Returns balance"
// @Failure 400 {object} map[string]string "Invalid account ID"
// @Failure 500 {object} map[string]string "Server error"
// @Router /accounts/{id}/balance [get]
func GetBalance(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	balance, err := repository.GetAccountBalance(database.DB, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
