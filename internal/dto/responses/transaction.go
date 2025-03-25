package responses

import "time"

type TransactionResponse struct {
    ID        int       `json:"id"`
    Amount    float64   `json:"amount"`
    Type      string    `json:"type"`
    Timestamp time.Time `json:"timestamp"`
    NewBalance float64  `json:"new_balance,omitempty"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}