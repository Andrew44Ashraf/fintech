package requests

type OpenAccountRequest struct {
    InitialBalance float64 `json:"initial_balance" validate:"gte=0"`
}