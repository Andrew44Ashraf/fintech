package handlers

import (
    "context"
    "github.com/Andrew44Ashraf/fintech-service/internal/dtos/responses"
    "github.com/Andrew44Ashraf/fintech-service/internal/repository"
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
    "strconv"
    "time"
)

type AccountHandler struct {
    accountRepo *repository.AccountRepository
}

func NewAccountHandler(repo *repository.AccountRepository) *AccountHandler {
    return &AccountHandler{accountRepo: repo}
}

// OpenAccount godoc
// @Summary Create a new account
// @Description Opens a new account with optional initial balance
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body requests.OpenAccountRequest false "Optional initial balance"
// @Success 200 {object} responses.AccountResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /accounts [post]
func (h *AccountHandler) OpenAccount(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()

    var req requests.OpenAccountRequest
    if err := c.ShouldBindJSON(&req); err != nil && err != io.EOF {
        c.JSON(http.StatusBadRequest, responses.NewErrorResponse("invalid request"))
        return
    }

    accountID, err := h.accountRepo.CreateAccount(ctx, req.InitialBalance)
    if err != nil {
        log.Printf("OpenAccount failed: %v", err)
        
        switch {
        case errors.Is(err, repository.ErrNegativeBalance):
            c.JSON(http.StatusBadRequest, responses.NewErrorResponse("initial balance cannot be negative"))
        default:
            c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("failed to create account"))
        }
        return
    }

    c.JSON(http.StatusOK, responses.AccountResponse{
        AccountID: accountID,
    })
}

// GetBalance godoc
// @Summary Get account balance
// @Tags accounts
// @Param id path int true "Account ID"
// @Success 200 {object} responses.BalanceResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /accounts/{id}/balance [get]
func (h *AccountHandler) GetBalance(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()

    accountID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, responses.NewErrorResponse("invalid account ID"))
        return
    }

    balance, err := h.accountRepo.GetAccountBalance(ctx, accountID)
    if err != nil {
        log.Printf("GetBalance failed: %v", err)
        
        switch {
        case errors.Is(err, repository.ErrAccountNotFound):
            c.JSON(http.StatusNotFound, responses.NewErrorResponse("account not found"))
        case errors.Is(err, repository.ErrAccountAlreadyClosed):
            c.JSON(http.StatusGone, responses.NewErrorResponse("account is closed"))
        default:
            c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("failed to get balance"))
        }
        return
    }

    c.JSON(http.StatusOK, responses.BalanceResponse{
        Balance: balance,
    })
}