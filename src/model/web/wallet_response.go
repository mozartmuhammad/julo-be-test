package web

import "time"

type WebResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type WalletCreateRequest struct {
	CustomerXID string `json:"name" validate:"required,min=1,max=36"`
}

type WalletResponse struct {
	ID        string     `json:"id"`
	OwnedBy   string     `json:"owned_by"`
	Status    string     `json:"status"`
	EnabledAt *time.Time `json:"enabled_at"`
	Balance   int        `json:"balance"`
}

type TransactionRequest struct {
	Amount      int    `json:"amount" validate:"required,min=1,numeric"`
	ReferenceID string `json:"reference_id" validate:"required,min=1"`
}

type TransactionResponse struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	TransactedAt time.Time `json:"transacted_at"`
	Type         string    `json:"type"`
	Amount       int       `json:"amount"`
	ReferenceID  string    `json:"reference_id"`
}

type DepositResponse struct {
	ID          string    `json:"id"`
	DepositedBy string    `json:"deposited_by"`
	Status      string    `json:"status"`
	DepositedAt time.Time `json:"deposited_at"`
	Amount      int       `json:"amount"`
	ReferenceID string    `json:"reference_id"`
}

type WithdrawalResponse struct {
	ID          string    `json:"id"`
	WithdrawnBy string    `json:"withdrawn_by"`
	Status      string    `json:"status"`
	WithdrawnAt time.Time `json:"withdrawn_at"`
	Amount      int       `json:"amount"`
	ReferenceID string    `json:"reference_id"`
}
