package repository

import (
	"context"
	"time"

	"github.com/mozartmuhammad/julo-be-test/model/domain"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet domain.Wallet) error
	GetWallet(ctx context.Context, customerXID string) (domain.Wallet, error)
	UpdateWalletStatus(ctx context.Context, customerXID string, status string, enabledAt *time.Time) error
	UpdateWalletBalance(ctx context.Context, walletID string, initialAmount, finalAmount int) (bool, error)

	GetWalletTransactions(ctx context.Context, walletID string) ([]domain.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, transactionID, status string) error
	AddTransaction(ctx context.Context, transaction domain.Transaction) error
}
