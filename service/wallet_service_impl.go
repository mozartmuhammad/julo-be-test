package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mozartmuhammad/julo-be-test/model/constants"
	"github.com/mozartmuhammad/julo-be-test/model/domain"
	"github.com/mozartmuhammad/julo-be-test/model/web"
	"github.com/mozartmuhammad/julo-be-test/repository"
)

type WalletService struct {
	WalletRepository repository.WalletRepository
	Validate         *validator.Validate
}

func NewWalletService(walletRepository repository.WalletRepository, DB *sql.DB, validate *validator.Validate) WalletServiceItf {
	return &WalletService{
		WalletRepository: walletRepository,
		Validate:         validate,
	}
}

func (svc *WalletService) InitializeWallet(ctx context.Context, request web.WalletCreateRequest) error {
	err := svc.Validate.Struct(request)
	if err != nil {
		return err
	}

	wallet := domain.Wallet{
		ID:          uuid.New().String(),
		CustomerXID: request.CustomerXID,
	}

	err = svc.WalletRepository.CreateWallet(ctx, wallet)
	if err != nil {
		return err
	}

	return nil
}

func (svc *WalletService) GetWalletBalance(ctx context.Context, customerXID string) (web.WalletResponse, error) {
	wallet, err := svc.WalletRepository.GetWallet(ctx, customerXID)
	if err != nil {
		return web.WalletResponse{}, err
	}

	return web.WalletResponse{
		ID:        wallet.ID,
		OwnedBy:   wallet.CustomerXID,
		Status:    wallet.Status,
		EnabledAt: wallet.EnabledAt,
		Balance:   wallet.Balance,
	}, nil
}

func (svc *WalletService) EnableWallet(ctx context.Context, customerXID string) (web.WalletResponse, error) {
	enabledAt := time.Now()
	err := svc.WalletRepository.UpdateWalletStatus(ctx, customerXID, constants.STATUS_ENABLED, &enabledAt)
	if err != nil {
		return web.WalletResponse{}, err
	}

	wallet, err := svc.WalletRepository.GetWallet(ctx, customerXID)
	if err != nil {
		return web.WalletResponse{}, err
	}

	return web.WalletResponse{
		ID:        wallet.ID,
		OwnedBy:   wallet.CustomerXID,
		Status:    wallet.Status,
		EnabledAt: wallet.EnabledAt,
		Balance:   wallet.Balance,
	}, nil
}

func (svc *WalletService) DisableWallet(ctx context.Context, customerXID string) (web.WalletResponse, error) {
	err := svc.WalletRepository.UpdateWalletStatus(ctx, customerXID, constants.STATUS_DISABLED, nil)
	if err != nil {
		return web.WalletResponse{}, err
	}

	wallet, err := svc.WalletRepository.GetWallet(ctx, customerXID)
	if err != nil {
		return web.WalletResponse{}, err
	}

	return web.WalletResponse{
		ID:        wallet.ID,
		OwnedBy:   wallet.CustomerXID,
		Status:    wallet.Status,
		EnabledAt: wallet.EnabledAt,
		Balance:   wallet.Balance,
	}, nil
}

func (svc *WalletService) GetWalletTransactions(ctx context.Context, customerXID string) ([]web.TransactionResponse, error) {
	wallet, err := svc.WalletRepository.GetWallet(ctx, customerXID)
	if err != nil {
		return []web.TransactionResponse{}, err
	}

	transaction, err := svc.WalletRepository.GetWalletTransactions(ctx, wallet.ID)
	if err != nil {
		return []web.TransactionResponse{}, err
	}

	result := []web.TransactionResponse{}
	for i := range transaction {
		result = append(result, web.TransactionResponse{
			ID:           transaction[i].ID,
			Status:       transaction[i].Status,
			TransactedAt: transaction[i].CreatedAt,
			Type:         transaction[i].TransactionType,
			Amount:       transaction[i].Amount,
			ReferenceID:  transaction[i].ReferenceID,
		})
	}
	return result, nil
}

func (svc *WalletService) AddWalletBalance(ctx context.Context, customerXID string, request web.TransactionRequest) (web.DepositResponse, error) {
	err := svc.Validate.StructCtx(ctx, request)
	if err != nil {
		return web.DepositResponse{}, err
	}

	wallet, err := svc.WalletRepository.GetWallet(ctx, customerXID)
	if err != nil {
		return web.DepositResponse{}, err
	}

	if wallet.Status == constants.STATUS_DISABLED {
		return web.DepositResponse{}, errors.New("wallet disabled")
	}

	isUpdated, err := svc.WalletRepository.UpdateWalletBalance(ctx, wallet.ID, wallet.Balance, wallet.Balance+request.Amount)
	if err != nil {
		return web.DepositResponse{}, err
	}
	status := constants.STATUS_FAILED
	if isUpdated {
		status = constants.STATUS_SUCCESS
	}

	transaction := domain.Transaction{
		ID:              uuid.New().String(),
		WalletID:        wallet.ID,
		CustomerXID:     wallet.CustomerXID,
		TransactionType: constants.TRANSACTION_TYPE_DEPOSIT,
		Amount:          request.Amount,
		ReferenceID:     request.ReferenceID,
		Status:          status,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err = svc.WalletRepository.AddTransaction(ctx, transaction)
	if err != nil {
		return web.DepositResponse{}, err
	}

	return web.DepositResponse{
		ID:          transaction.ID,
		DepositedBy: transaction.CustomerXID,
		Status:      transaction.Status,
		DepositedAt: transaction.CreatedAt,
		Amount:      transaction.Amount,
		ReferenceID: transaction.ReferenceID,
	}, nil
}

func (svc *WalletService) DeductWalletBalance(ctx context.Context, customerXID string, request web.TransactionRequest) (web.WithdrawalResponse, error) {
	wallet, err := svc.WalletRepository.GetWallet(ctx, customerXID)
	if err != nil {
		return web.WithdrawalResponse{}, err
	}

	if wallet.Status == constants.STATUS_DISABLED {
		return web.WithdrawalResponse{}, errors.New("wallet disabled")
	}

	if request.Amount > wallet.Balance {
		return web.WithdrawalResponse{}, errors.New("insufficient balance")
	}

	isUpdated, err := svc.WalletRepository.UpdateWalletBalance(ctx, wallet.ID, wallet.Balance, wallet.Balance-request.Amount)
	if err != nil {
		return web.WithdrawalResponse{}, err
	}
	status := constants.STATUS_FAILED
	if isUpdated {
		status = constants.STATUS_SUCCESS
	}

	transaction := domain.Transaction{
		ID:              uuid.New().String(),
		WalletID:        wallet.ID,
		CustomerXID:     wallet.CustomerXID,
		TransactionType: constants.TRANSACTION_TYPE_WITHDRAWAL,
		Amount:          request.Amount,
		ReferenceID:     request.ReferenceID,
		Status:          status,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err = svc.WalletRepository.AddTransaction(ctx, transaction)
	if err != nil {
		return web.WithdrawalResponse{}, err
	}

	return web.WithdrawalResponse{
		ID:          transaction.ID,
		WithdrawnBy: transaction.CustomerXID,
		Status:      transaction.Status,
		WithdrawnAt: transaction.CreatedAt,
		Amount:      transaction.Amount,
		ReferenceID: transaction.ReferenceID,
	}, nil
}
