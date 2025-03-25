package repository

import (
	"database/sql"
	"fintech-service/internal/models"
)

func CreateAccount(db *sql.DB) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO accounts (balance) VALUES (0) RETURNING id").Scan(&id)
	return id, err
}

func GetAccountBalance(db *sql.DB, accountID int) (float64, error) {
	var balance float64
	err := db.QueryRow("SELECT balance FROM accounts WHERE id = $1", accountID).Scan(&balance)
	return balance, err
}
