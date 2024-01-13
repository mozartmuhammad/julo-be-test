package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mozartmuhammad/julo-be-test/src/model/constants"
	"github.com/mozartmuhammad/julo-be-test/src/model/domain"
	"github.com/mozartmuhammad/julo-be-test/src/model/web"
	"github.com/mozartmuhammad/julo-be-test/src/repository"
)

type WalletService struct {
	WalletRepository repository.WalletRepository
	Validate         *validator.Validate
	DelayDuration    time.Duration
}

func NewWalletService(walletRepository repository.WalletRepository, validate *validator.Validate, delay time.Duration) WalletServiceItf {
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

	// create new wallet
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
	wallet, err := svc.WalletRepository.GetWallet(ctx, customerXID)
	if err != nil {
		return web.WalletResponse{}, err
	}

	// check wallet status
	if wallet.Status == constants.STATUS_ENABLED {
		return web.WalletResponse{}, errors.New("Already enabled")
	}

	// update wallet status
	enabledAt := time.Now()
	err = svc.WalletRepository.UpdateWalletStatus(ctx, customerXID, constants.STATUS_ENABLED, &enabledAt)
	if err != nil {
		return web.WalletResponse{}, err
	}

	// get latest wallet data
	wallet, err = svc.WalletRepository.GetWallet(ctx, customerXID)
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
	// disable wallet
	err := svc.WalletRepository.UpdateWalletStatus(ctx, customerXID, constants.STATUS_DISABLED, nil)
	if err != nil {
		return web.WalletResponse{}, err
	}

	// get wallet data
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

	// check wallet status
	if wallet.Status == constants.STATUS_DISABLED {
		return []web.TransactionResponse{}, errors.New("wallet disabled")
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

	// check wallet status
	if wallet.Status == constants.STATUS_DISABLED {
		return web.DepositResponse{}, errors.New("wallet disabled")
	}

	// insert transaction with status pending
	transaction := domain.Transaction{
		ID:              uuid.New().String(),
		WalletID:        wallet.ID,
		CustomerXID:     wallet.CustomerXID,
		TransactionType: constants.TRANSACTION_TYPE_DEPOSIT,
		Amount:          request.Amount,
		ReferenceID:     request.ReferenceID,
		Status:          constants.STATUS_PENDING,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err = svc.WalletRepository.AddTransaction(ctx, transaction)
	if err != nil {
		return web.DepositResponse{}, err
	}

	go func() {
		// delay 5 seconds for update wallet balance
		time.Sleep(svc.DelayDuration)
		isUpdated, err := svc.WalletRepository.UpdateWalletBalance(context.Background(), wallet.ID, wallet.Balance, wallet.Balance+request.Amount)
		if err != nil {
			log.Println("error update wallet balance:", err.Error())
		}

		status := constants.STATUS_FAILED
		if isUpdated {
			status = constants.STATUS_SUCCESS
		}

		err = svc.WalletRepository.UpdateTransactionStatus(context.Background(), transaction.ID, status)
		if err != nil {
			log.Println("error update transaction status:", err.Error())
		}
	}()

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
	err := svc.Validate.StructCtx(ctx, request)
	if err != nil {
		return web.WithdrawalResponse{}, err
	}

	wallet, err := svc.WalletRepository.GetWallet(ctx, customerXID)
	if err != nil {
		return web.WithdrawalResponse{}, err
	}

	// check wallet status
	if wallet.Status == constants.STATUS_DISABLED {
		return web.WithdrawalResponse{}, errors.New("wallet disabled")
	}

	// compare amount with balance
	if request.Amount > wallet.Balance {
		return web.WithdrawalResponse{}, errors.New("insufficient balance")
	}

	// insert transaction with status pending
	transaction := domain.Transaction{
		ID:              uuid.New().String(),
		WalletID:        wallet.ID,
		CustomerXID:     wallet.CustomerXID,
		TransactionType: constants.TRANSACTION_TYPE_WITHDRAWAL,
		Amount:          request.Amount,
		ReferenceID:     request.ReferenceID,
		Status:          constants.STATUS_PENDING,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err = svc.WalletRepository.AddTransaction(ctx, transaction)
	if err != nil {
		return web.WithdrawalResponse{}, err
	}

	go func() {
		time.Sleep(svc.DelayDuration)
		isUpdated, err := svc.WalletRepository.UpdateWalletBalance(context.Background(), wallet.ID, wallet.Balance, wallet.Balance-request.Amount)
		if err != nil {
			log.Println("error update wallet balance:", err.Error())
		}

		status := constants.STATUS_FAILED
		if isUpdated {
			status = constants.STATUS_SUCCESS
		}
		err = svc.WalletRepository.UpdateTransactionStatus(context.Background(), transaction.ID, status)
		if err != nil {
			log.Println("error update transaction status:", err.Error())
		}
	}()

	return web.WithdrawalResponse{
		ID:          transaction.ID,
		WithdrawnBy: transaction.CustomerXID,
		Status:      transaction.Status,
		WithdrawnAt: transaction.CreatedAt,
		Amount:      transaction.Amount,
		ReferenceID: transaction.ReferenceID,
	}, nil
}
