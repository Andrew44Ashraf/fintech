package repository

import (
	"database/sql"
	"errors"
	"fintech-service/internal/models"
	"log"
)

// CreateTransaction ensures funds are updated atomically using SQL transactions
func CreateTransaction(db *sql.DB, accountID int, amount float64, txType models.TransactionType) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // Ensures rollback on failure

	var balance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", accountID).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("account not found")
		}
		return 0, err
	}

	if txType == models.Withdrawal && balance < amount {
		return 0, errors.New("insufficient funds")
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, accountID)
	if err != nil {
		return 0, err
	}

	var txID int
	err = tx.QueryRow(
		"INSERT INTO transactions (account_id, amount, type) VALUES ($1, $2, $3) RETURNING id",
		accountID, amount, txType,
	).Scan(&txID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	log.Printf("Transaction successful: %s of %.2f for account %d", txType, amount, accountID)
	return txID, nil
}

// GetTransactions supports pagination
func GetTransactions(db *sql.DB, accountID int, limit int, offset int) ([]models.Transaction, error) {
	rows, err := db.Query("SELECT id, account_id, amount, type, timestamp FROM transactions WHERE account_id = $1 ORDER BY timestamp DESC LIMIT $2 OFFSET $3",
		accountID, limit, offset)
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
