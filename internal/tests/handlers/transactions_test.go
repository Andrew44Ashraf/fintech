package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/Andrew44Ashraf/fintech-service/internal/handlers"
	"github.com/Andrew44Ashraf/fintech-service/internal/repository"
	"github.com/Andrew44Ashraf/fintech-service/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestDepositHandler(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	accountRepo, _ := testutils.NewMockRepository()
	transactionRepo, mock := testutils.NewMockTransactionRepository()
	handler := handlers.NewTransactionHandler(transactionRepo, accountRepo)

	t.Run("successful deposit", func(t *testing.T) {
		// Mock expectations
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO transactions`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectExec(`UPDATE accounts`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Create test request
		payload := `{"amount": 100.0}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/accounts/1/deposit", bytes.NewBufferString(payload))
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		// Call handler
		handler.Deposit(c)

		// Verify
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
		
		var response handlers.TransactionResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 1, response.TransactionID)
	})

	t.Run("invalid amount", func(t *testing.T) {
		payload := `{"amount": -50.0}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/accounts/1/deposit", bytes.NewBufferString(payload))
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		handler.Deposit(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}