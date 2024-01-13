package domain

import "time"

type Wallet struct {
	ID          string
	CustomerXID string
	Status      string
	EnabledAt   *time.Time
	Balance     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Transaction struct {
	ID              string
	WalletID        string
	CustomerXID     string
	TransactionType string
	Amount          int
	ReferenceID     string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
