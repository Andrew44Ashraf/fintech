package handlers

import (
	"context"
	"github.com/Andrew44Ashraf/fintech-service/internal/dtos/requests"
	"github.com/Andrew44Ashraf/fintech-service/internal/dtos/responses"
	"github.com/Andrew44Ashraf/fintech-service/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TransactionHandler struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
}

func NewTransactionHandler(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
) *TransactionHandler {
	return &TransactionHandler{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

// Deposit godoc
// @Summary Deposit funds
// @Description Deposits money into an account
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param body body requests.DepositRequest true "Deposit amount"
// @Success 200 {object} responses.TransactionResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /accounts/{id}/deposit [post]
func (h *TransactionHandler) Deposit(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("invalid account ID"))
		return
	}

	var req requests.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Deposit: invalid input - %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("invalid request body"))
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error()))
		return
	}

	// Process deposit
	txID, err := h.transactionRepo.CreateDeposit(ctx, accountID, req.Amount)
	if err != nil {
		log.Printf("Deposit failed: %v", err)
		
		switch {
		case errors.Is(err, repository.ErrAccountNotFound):
			c.JSON(http.StatusNotFound, responses.NewErrorResponse("account not found"))
		case errors.Is(err, repository.ErrAccountClosed):
			c.JSON(http.StatusForbidden, responses.NewErrorResponse("account is closed"))
		default:
			c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("failed to process deposit"))
		}
		return
	}

	// Get updated balance
	balance, err := h.accountRepo.GetAccountBalance(ctx, accountID)
	if err != nil {
		log.Printf("Deposit: failed to get updated balance: %v", err)
		balance = 0 // Continue response without balance
	}

	c.JSON(http.StatusOK, responses.TransactionResponse{
		TransactionID: txID,
		NewBalance:    balance,
		Message:       "Deposit processed successfully",
	})
}

// Withdraw godoc
// @Summary Withdraw funds
// @Description Withdraws money from an account
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param body body requests.WithdrawRequest true "Withdrawal amount"
// @Success 200 {object} responses.TransactionResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /accounts/{id}/withdraw [post]
func (h *TransactionHandler) Withdraw(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("invalid account ID"))
		return
	}

	var req requests.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Withdraw: invalid input - %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("invalid request body"))
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error()))
		return
	}

	// Process withdrawal
	txID, err := h.transactionRepo.CreateWithdrawal(ctx, accountID, req.Amount)
	if err != nil {
		log.Printf("Withdraw failed: %v", err)
		
		switch {
		case errors.Is(err, repository.ErrAccountNotFound):
			c.JSON(http.StatusNotFound, responses.NewErrorResponse("account not found"))
		case errors.Is(err, repository.ErrAccountClosed):
			c.JSON(http.StatusForbidden, responses.NewErrorResponse("account is closed"))
		case errors.Is(err, repository.ErrInsufficientFunds):
			c.JSON(http.StatusBadRequest, responses.NewErrorResponse("insufficient funds"))
		default:
			c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("failed to process withdrawal"))
		}
		return
	}

	// Get updated balance
	balance, err := h.accountRepo.GetAccountBalance(ctx, accountID)
	if err != nil {
		log.Printf("Withdraw: failed to get updated balance: %v", err)
		balance = 0 // Continue response without balance
	}

	c.JSON(http.StatusOK, responses.TransactionResponse{
		TransactionID: txID,
		NewBalance:    balance,
		Message:       "Withdrawal processed successfully",
	})
}