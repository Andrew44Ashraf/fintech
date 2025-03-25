package repository

import (
	"database/sql"
	"errors"
	"github.com/Andrew44Ashraf/fintech-service/internal/errors"
	"github.com/Andrew44Ashraf/fintech-service/internal/models"
	"fmt"
	"time"
)

// CreateTransaction ensures funds are updated atomically using SQL transactions
func CreateTransaction(db *sql.DB, accountID int, amount float64, txType models.TransactionType) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() // Ensures rollback on failure

	var balance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", accountID).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("account ID %d: %w", accountID, errors.ErrAccountNotFound)
		}
		return 0, fmt.Errorf("failed to fetch account balance: %w", err)
	}

	// Ensure sufficient funds for withdrawal
	if txType == models.Withdrawal && balance < amount {
		return 0, fmt.Errorf("withdrawal of %.2f failed: %w", amount, errors.ErrInsufficientFunds)
	}

	// Update balance
	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, accountID)
	if err != nil {
		return 0, fmt.Errorf("failed to update balance for account %d: %w", accountID, err)
	}

	// Insert transaction
	var txID int
	err = tx.QueryRow("INSERT INTO transactions (account_id, amount, type) VALUES ($1, $2, $3) RETURNING id",
		accountID, amount, txType).Scan(&txID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert transaction: %w", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return txID, nil
}

// GetTransactions retrieves transactions for an account with pagination
func GetTransactions(db *sql.DB, accountID int, limit int, offset int) ([]models.Transaction, error) {
	rows, err := db.Query(`
		SELECT id, account_id, amount, type, timestamp 
		FROM transactions 
		WHERE account_id = $1 
		ORDER BY timestamp DESC 
		LIMIT $2 OFFSET $3`, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.AccountID, &tx.Amount, &tx.Type, &tx.Timestamp); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

// GetTransactionsByDateRange retrieves transactions within a specific date range
func GetTransactionsByDateRange(db *sql.DB, accountID int, startDate, endDate time.Time) ([]models.Transaction, error) {
	rows, err := db.Query(`
		SELECT id, account_id, amount, type, timestamp 
		FROM transactions 
		WHERE account_id = $1 AND timestamp BETWEEN $2 AND $3 
		ORDER BY timestamp DESC`, accountID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.AccountID, &tx.Amount, &tx.Type, &tx.Timestamp); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

// GetTotalDepositsAndWithdrawals returns the total deposits and withdrawals for an account
func GetTotalDepositsAndWithdrawals(db *sql.DB, accountID int) (float64, float64, error) {
	var totalDeposits, totalWithdrawals float64

	err := db.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'deposit' THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'withdrawal' THEN amount ELSE 0 END), 0)
		FROM transactions
		WHERE account_id = $1`, accountID).Scan(&totalDeposits, &totalWithdrawals)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to calculate totals: %w", err)
	}

	return totalDeposits, totalWithdrawals, nil
}
