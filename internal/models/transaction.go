package models

import (
	"time"
)

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
)

// Transaction represents the database model
type Transaction struct {
	ID        int             `json:"id" gorm:"primaryKey"`
	AccountID int             `json:"account_id" validate:"required" gorm:"index"`
	Amount    float64         `json:"amount" validate:"required,gt=0" gorm:"type:decimal(15,2)"`
	Type      TransactionType `json:"type" validate:"required,oneof=deposit withdrawal" gorm:"type:varchar(20)"`
	Timestamp time.Time       `json:"timestamp" gorm:"autoCreateTime"`
}

// DepositRequest defines the payload for deposit operations
type DepositRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

// WithdrawRequest defines the payload for withdrawal operations
type WithdrawRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

// TransactionResponse defines the API response format
type TransactionResponse struct {
	ID        int             `json:"id"`
	AccountID int             `json:"account_id,omitempty"` // Optional in responses
	Amount    float64         `json:"amount"`
	Type      TransactionType `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	NewBalance float64        `json:"new_balance,omitempty"` // Computed field
	Message   string          `json:"message,omitempty"`     // Optional status message
}

// TransactionListResponse for paginated results
type TransactionListResponse struct {
	Data         []TransactionResponse `json:"data"`
	TotalRecords int                   `json:"total_records"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
}

// Helper function to convert DB model to API response
func (t *Transaction) ToResponse(newBalance float64) TransactionResponse {
	return TransactionResponse{
		ID:         t.ID,
		AccountID:  t.AccountID,
		Amount:     t.Amount,
		Type:       t.Type,
		Timestamp:  t.Timestamp,
		NewBalance: newBalance,
	}
}