package responses

type AccountResponse struct {
    AccountID int `json:"account_id"`
}

type BalanceResponse struct {
    Balance float64 `json:"balance"`
}