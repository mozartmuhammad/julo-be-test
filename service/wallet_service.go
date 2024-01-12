package service

import (
	"context"

	"github.com/mozartmuhammad/julo-be-test/model/web"
)

type WalletServiceItf interface {
	InitializeWallet(ctx context.Context, request web.WalletCreateRequest) error
	GetWalletBalance(ctx context.Context, customerXID string) (web.WalletResponse, error)
	EnableWallet(ctx context.Context, customerXID string) (web.WalletResponse, error)
	DisableWallet(ctx context.Context, customerXID string) (web.WalletResponse, error)
	GetWalletTransactions(ctx context.Context, customerXID string) ([]web.TransactionResponse, error)
	AddWalletBalance(ctx context.Context, customerXID string, request web.TransactionRequest) (web.DepositResponse, error)
	DeductWalletBalance(ctx context.Context, customerXID string, request web.TransactionRequest) (web.WithdrawalResponse, error)
}
