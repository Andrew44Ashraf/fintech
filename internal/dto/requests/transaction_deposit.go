package requests

import "github.com/go-playground/validator/v10"

type DepositRequest struct {
    Amount float64 `json:"amount" validate:"required,gt=0"`
}

func (r *DepositRequest) Validate() error {
    validate := validator.New()
    return validate.Struct(r)
}