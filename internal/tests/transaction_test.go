package tests

import (
	"fintech-service/internal/models"
	"fintech-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeposit(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	// Create a test account
	accountID, _ := repository.CreateAccount(db)

	// Perform a deposit
	txID, err := repository.CreateTransaction(db, accountID, 100.0, models.Deposit)
	assert.NoError(t, err, "Expected deposit to succeed")
	assert.NotZero(t, txID, "Transaction ID should not be zero")

	// Verify the balance
	balance, err := repository.GetAccountBalance(db, accountID)
	assert.NoError(t, err)
	assert.Equal(t, 100.0, balance, "Balance should be updated after deposit")
}

func TestWithdrawWithInsufficientFunds(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	accountID, _ := repository.CreateAccount(db)

	// Attempt withdrawal of more than balance
	_, err := repository.CreateTransaction(db, accountID, 500.0, models.Withdrawal)
	assert.Error(t, err, "Expected error when withdrawing more than balance")
}
