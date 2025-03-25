package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrTransactionFailed   = errors.New("transaction processing failed")
	ErrNegativeAmount      = errors.New("amount must be positive")
	ErrInvalidTransaction = errors.New("invalid transaction type")
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Transaction represents a financial transaction
type Transaction struct {
	ID          int
	AccountID   int
	Amount      float64
	Type        string // "deposit" or "withdrawal"
	CreatedAt   time.Time
	FinalBalance float64
}

// CreateDeposit handles deposit transactions atomically
func (r *TransactionRepository) CreateDeposit(ctx context.Context, accountID int, amount float64) (int, error) {
	if amount <= 0 {
		return 0, ErrNegativeAmount
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Verify account exists and is active
	var accountStatus string
	err = tx.QueryRowContext(ctx,
		"SELECT status FROM accounts WHERE id = $1 FOR UPDATE",
		accountID,
	).Scan(&accountStatus)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return 0, ErrAccountNotFound
	case err != nil:
		return 0, fmt.Errorf("account verification failed: %w", err)
	case accountStatus != "active":
		return 0, ErrAccountClosed
	}

	// 2. Create transaction record
	var txID int
	err = tx.QueryRowContext(ctx,
		`INSERT INTO transactions 
		 (account_id, amount, type) 
		 VALUES ($1, $2, 'deposit')
		 RETURNING id`,
		accountID, amount,
	).Scan(&txID)
	if err != nil {
		return 0, fmt.Errorf("failed to create transaction: %w", err)
	}

	// 3. Update account balance
	_, err = tx.ExecContext(ctx,
		"UPDATE accounts SET balance = balance + $1 WHERE id = $2",
		amount, accountID,
	)
	if err != nil {
		return 0, fmt.Errorf("balance update failed: %w", err)
	}

	// 4. Get final balance for the response
	var finalBalance float64
	err = tx.QueryRowContext(ctx,
		"SELECT balance FROM accounts WHERE id = $1",
		accountID,
	).Scan(&finalBalance)
	if err != nil {
		return 0, fmt.Errorf("failed to get final balance: %w", err)
	}

	// 5. Update transaction with final balance
	_, err = tx.ExecContext(ctx,
		"UPDATE transactions SET final_balance = $1 WHERE id = $2",
		finalBalance, txID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to update transaction record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("transaction commit failed: %w", err)
	}

	return txID, nil
}

// CreateWithdrawal handles withdrawal transactions atomically
func (r *TransactionRepository) CreateWithdrawal(ctx context.Context, accountID int, amount float64) (int, error) {
	if amount <= 0 {
		return 0, ErrNegativeAmount
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Verify account status and get current balance
	var (
		currentBalance float64
		accountStatus  string
	)
	err = tx.QueryRowContext(ctx,
		"SELECT balance, status FROM accounts WHERE id = $1 FOR UPDATE",
		accountID,
	).Scan(&currentBalance, &accountStatus)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return 0, ErrAccountNotFound
	case err != nil:
		return 0, fmt.Errorf("account verification failed: %w", err)
	case accountStatus != "active":
		return 0, ErrAccountClosed
	case currentBalance < amount:
		return 0, ErrInsufficientFunds
	}

	// 2. Create transaction record
	var txID int
	err = tx.QueryRowContext(ctx,
		`INSERT INTO transactions 
		 (account_id, amount, type) 
		 VALUES ($1, $2, 'withdrawal')
		 RETURNING id`,
		accountID, amount,
	).Scan(&txID)
	if err != nil {
		return 0, fmt.Errorf("failed to create transaction: %w", err)
	}

	// 3. Update account balance
	_, err = tx.ExecContext(ctx,
		"UPDATE accounts SET balance = balance - $1 WHERE id = $2",
		amount, accountID,
	)
	if err != nil {
		return 0, fmt.Errorf("balance update failed: %w", err)
	}

	// 4. Get final balance
	var finalBalance float64
	err = tx.QueryRowContext(ctx,
		"SELECT balance FROM accounts WHERE id = $1",
		accountID,
	).Scan(&finalBalance)
	if err != nil {
		return 0, fmt.Errorf("failed to get final balance: %w", err)
	}

	// 5. Update transaction record
	_, err = tx.ExecContext(ctx,
		"UPDATE transactions SET final_balance = $1 WHERE id = $2",
		finalBalance, txID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to update transaction record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("transaction commit failed: %w", err)
	}

	return txID, nil
}

// GetTransactions retrieves transaction history for an account
func (r *TransactionRepository) GetTransactions(ctx context.Context, accountID int, limit, offset int) ([]Transaction, error) {
	const query = `
		SELECT id, account_id, amount, type, created_at, final_balance
		FROM transactions
		WHERE account_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(
			&t.ID,
			&t.AccountID,
			&t.Amount,
			&t.Type,
			&t.CreatedAt,
			&t.FinalBalance,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return transactions, nil
}