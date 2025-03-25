package models

import "time"

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
)

type Transaction struct {
	ID        int             `json:"id"`
	AccountID int             `json:"account_id" validate:"required"`
	Amount    float64         `json:"amount" validate:"required,gt=0"`
	Type      TransactionType `json:"type" validate:"required,oneof=deposit withdrawal"`
	Timestamp time.Time       `json:"timestamp"`
}
