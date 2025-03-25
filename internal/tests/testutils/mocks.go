package testutils

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Andrew44Ashraf/fintech-service/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewMockDB creates a sqlmock database connection
func NewMockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("failed to create sqlmock")
	}
	return db, mock
}

// NewMockRepository creates a repository with mock DB
func NewMockRepository() (*repository.AccountRepository, sqlmock.Sqlmock) {
	db, mock := NewMockDB()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		panic("failed to open gorm db")
	}
	return repository.NewAccountRepository(gormDB), mock
}