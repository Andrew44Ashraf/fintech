package handlers

import (
	"fintech-service/internal/database"
	"fintech-service/internal/models"
	"fintech-service/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
)

// Deposit godoc
// @Summary Deposit funds
// @Description Deposits money into an account
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param body body struct{Amount float64} true "Deposit amount"
// @Success 200 {object} map[string]int "Transaction ID"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Server error"
// @Router /accounts/{id}/deposit [post]
func Deposit(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Invalid deposit input:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	// Validate using go-playground/validator
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txID, err := repository.CreateTransaction(database.DB, accountID, input.Amount, models.Deposit)
	if err != nil {
		log.Println("Deposit error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction_id": txID})
}

// Withdraw godoc
// @Summary Withdraw funds
// @Description Withdraws money from an account
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param body body struct{Amount float64} true "Withdrawal amount"
// @Success 200 {object} map[string]int "Transaction ID"
// @Failure 400 {object} map[string]string "Invalid input or insufficient funds"
// @Failure 500 {object} map[string]string "Server error"
// @Router /accounts/{id}/withdraw [post]
func Withdraw(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	// Validate using go-playground/validator
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := repository.GetAccountBalance(database.DB, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "account not found"})
		return
	}

	if balance < input.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
		return
	}

	txID, err := repository.CreateTransaction(database.DB, accountID, -input.Amount, models.Withdrawal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction_id": txID})
}
