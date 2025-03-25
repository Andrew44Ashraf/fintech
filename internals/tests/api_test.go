package tests

import (
	"bytes"
	"fintech-service/internal/database"
	"fintech-service/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	routes.SetupRoutes(r)
	return r
}

func TestOpenAccountAPI(t *testing.T) {
	database.InitDB()
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/api/accounts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDepositAPI(t *testing.T) {
	database.InitDB()
	router := setupRouter()

	// Create account first
	accountID := 1

	// Perform a deposit
	body := []byte(`{"amount": 100}`)
	req, _ := http.NewRequest("POST", "/api/accounts/1/deposit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
