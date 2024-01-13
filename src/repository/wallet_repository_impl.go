package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/mozartmuhammad/julo-be-test/src/model/domain"
)

type WalletRepositoryImpl struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) WalletRepository {
	return &WalletRepositoryImpl{
		db: db,
	}
}

func (repo *WalletRepositoryImpl) CreateWallet(ctx context.Context, wallet domain.Wallet) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, insertWalletQuery, wallet.ID, wallet.CustomerXID)
	if err != nil {
		return err
	}

	errorCommit := tx.Commit()
	if errorCommit != nil {
		_ = tx.Rollback()

		return errorCommit
	}

	return nil
}

func (repo *WalletRepositoryImpl) UpdateWalletStatus(ctx context.Context, customerXID string, status string, enabledAt *time.Time) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, updateWalletStatusQuery, status, enabledAt, customerXID)
	if err != nil {
		return err
	}

	errorCommit := tx.Commit()
	if errorCommit != nil {
		_ = tx.Rollback()

		return errorCommit
	}

	return nil
}

func (repo *WalletRepositoryImpl) GetWallet(ctx context.Context, customerXID string) (domain.Wallet, error) {
	var result domain.Wallet
	err := repo.db.QueryRowContext(ctx, getWalletQuery, customerXID).Scan(
		&result.ID,
		&result.CustomerXID,
		&result.Status,
		&result.EnabledAt,
		&result.Balance,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repo *WalletRepositoryImpl) GetWalletTransactions(ctx context.Context, walletID string) ([]domain.Transaction, error) {
	var result []domain.Transaction
	rows, err := repo.db.QueryContext(ctx, getTransactionsQuery, walletID)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		data := domain.Transaction{}
		err := rows.Scan(
			&data.ID,
			&data.WalletID,
			&data.CustomerXID,
			&data.TransactionType,
			&data.Amount,
			&data.ReferenceID,
			&data.Status,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			return result, err
		}
		result = append(result, data)
	}
	return result, nil
}

func (repo *WalletRepositoryImpl) AddTransaction(ctx context.Context, transaction domain.Transaction) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, insertTransactionQuery,
		transaction.ID,
		transaction.WalletID,
		transaction.CustomerXID,
		transaction.TransactionType,
		transaction.Amount,
		transaction.ReferenceID,
		transaction.Status,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)
	if err != nil {
		return err
	}

	errorCommit := tx.Commit()
	if errorCommit != nil {
		_ = tx.Rollback()

		return errorCommit
	}

	return nil
}

func (repo *WalletRepositoryImpl) UpdateTransactionStatus(ctx context.Context, transactionID string, status string) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, updateTransactionStatusQuery, status, transactionID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repo *WalletRepositoryImpl) UpdateWalletBalance(ctx context.Context, walletID string, initialAmount, finalAmount int) (bool, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return false, err
	}

	sss, err := tx.ExecContext(ctx, updateWalletBalanceQuery, finalAmount, walletID, initialAmount)
	if err != nil {
		return false, err
	}

	rowsAffected, _ := sss.RowsAffected()
	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
