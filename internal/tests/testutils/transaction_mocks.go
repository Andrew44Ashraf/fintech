package testutils

import (
	"database/sql"
	"github.com/Andrew44Ashraf/fintech-service/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
)

func NewMockTransactionRepository() (*repository.TransactionRepository, sqlmock.Sqlmock) {
	db, mock := NewMockDB()
	return repository.NewTransactionRepository(db), mock
}