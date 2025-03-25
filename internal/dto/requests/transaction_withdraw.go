package requests

type WithdrawRequest struct {
    Amount float64 `json:"amount" validate:"required,gt=0"`
}