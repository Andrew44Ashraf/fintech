package tests

import (
	"database/sql"
	"github.com/Andrew44Ashraf/fintech-service/internal/database"
	"github.com/Andrew44Ashraf/fintech-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func setupTestDB() *sql.DB {
	// Use an in-memory database for testing
	db, err := sql.Open("postgres", "host=localhost user=postgres password=password dbname=fintech_db sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}
	return db
}

func TestCreateAccount(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	accountID, err := repository.CreateAccount(db)
	assert.NoError(t, err, "Expected no error when creating account")
	assert.NotZero(t, accountID, "Account ID should be greater than 0")
}
