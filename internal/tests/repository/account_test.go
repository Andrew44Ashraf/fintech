package repository_test

import (
	"testing"
	"time"
	"github.com/Andrew44Ashraf/fintech-service/internal/repository"
	"github.com/Andrew44Ashraf/fintech-service/internal/testutils"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	repo, mock := testutils.NewMockRepository()

	t.Run("successful account creation", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO accounts`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		id, err := repo.CreateAccount(context.Background(), 100.0)
		
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("negative balance should fail", func(t *testing.T) {
		_, err := repo.CreateAccount(context.Background(), -50.0)
		assert.ErrorIs(t, err, repository.ErrNegativeBalance)
	})
}

func TestGetAccountBalance(t *testing.T) {
	repo, mock := testutils.NewMockRepository()

	t.Run("existing account", func(t *testing.T) {
		mock.ExpectQuery(`SELECT balance FROM accounts`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(100.0))

		balance, err := repo.GetAccountBalance(context.Background(), 1)
		
		assert.NoError(t, err)
		assert.Equal(t, 100.0, balance)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("non-existent account", func(t *testing.T) {
		mock.ExpectQuery(`SELECT balance FROM accounts`).
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetAccountBalance(context.Background(), 999)
		
		assert.ErrorIs(t, err, repository.ErrAccountNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}