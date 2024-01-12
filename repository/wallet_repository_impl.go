package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mozartmuhammad/julo-be-test/model/domain"
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
	fmt.Println("haha")
	tx, err := repo.db.Begin()
	if err != nil {
		fmt.Println("hihi")
		return err
	}

	SQL := `INSERT INTO wallets
		(id, customer_xid)
		VALUES(?, ?)`

	_, err = tx.ExecContext(ctx, SQL, wallet.ID, wallet.CustomerXID)
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

	SQL := `UPDATE wallets
		SET
			status = ?,
			enabled_at = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE 
			customer_xid = ?`

	_, err = tx.ExecContext(ctx, SQL, status, enabledAt, customerXID)
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
	SQL := "select id, customer_xid, status, enabled_at, balance, created_at, updated_at FROM wallets WHERE customer_xid = ?"
	err := repo.db.QueryRowContext(ctx, SQL, customerXID).Scan(
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
	SQL := "select id, wallet_id, customer_xid, transaction_type, amount, reference_id, status, created_at, updated_at FROM transactions WHERE wallet_id = ? order by created_at"
	rows, err := repo.db.QueryContext(ctx, SQL, walletID)
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
			&data.Status,
			&data.ReferenceID,
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

	SQL := `INSERT INTO transactions
		(id, wallet_id, customer_xid, transaction_type, amount, reference_id, status, created_at, updated_at)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, SQL,
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
