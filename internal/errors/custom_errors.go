package errors

import "errors"

// Predefined errors for standardization
var (
	ErrAccountNotFound    = errors.New("account not found")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrInvalidTransaction = errors.New("invalid transaction")
)
