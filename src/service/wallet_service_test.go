package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock_repository "github.com/mozartmuhammad/julo-be-test/src/mock/repository"
	"github.com/mozartmuhammad/julo-be-test/src/model/domain"
	"github.com/mozartmuhammad/julo-be-test/src/model/web"
	"github.com/mozartmuhammad/julo-be-test/src/service"
)

var (
	svc service.WalletServiceItf

	mockRepository *mock_repository.MockWalletRepository
)

func provideTest(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository = mock_repository.NewMockWalletRepository(ctrl)
	validator := validator.New()
	svc = service.NewWalletService(mockRepository, validator, 0)

	return func() {}
}

func TestInitializeWallet(t *testing.T) {
	type (
		args struct {
			payload web.WalletCreateRequest
		}
	)

	testCases := []struct {
		testID   int
		testDesc string
		args     args
		mockFunc func()
		wantErr  bool
	}{
		{
			testID:   1,
			testDesc: "Success",
			args: args{
				payload: web.WalletCreateRequest{
					CustomerXID: "abcdef",
				},
			},
			mockFunc: func() {
				mockRepository.EXPECT().CreateWallet(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			testID:   2,
			testDesc: "Failed - error call repo CreateWallet",
			args: args{
				payload: web.WalletCreateRequest{
					CustomerXID: "abcdef",
				},
			},
			mockFunc: func() {
				mockRepository.EXPECT().CreateWallet(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			testID:   3,
			testDesc: "Failed - error validate",
			args: args{
				payload: web.WalletCreateRequest{
					CustomerXID: "",
				},
			},
			mockFunc: func() {
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDesc, func(t *testing.T) {
			testDep := provideTest(t)
			defer testDep()
			tc.mockFunc()

			err := svc.InitializeWallet(context.Background(), tc.args.payload)
			assert.Equal(t, err != nil, tc.wantErr)
		})
	}
}

func TestGetWalletBalance(t *testing.T) {
	type (
		args struct {
			customerXID string
		}
	)

	testCases := []struct {
		testID     int
		testDesc   string
		args       args
		mockFunc   func()
		wantErr    bool
		wantResult web.WalletResponse
	}{
		{
			testID:   1,
			testDesc: "Success",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID: "mock-id",
				}, nil)
			},
			wantErr: false,
			wantResult: web.WalletResponse{
				ID: "mock-id",
			},
		},
		{
			testID:   2,
			testDesc: "Failed - error call GetWallet",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{}, fmt.Errorf("error"))
			},
			wantErr:    true,
			wantResult: web.WalletResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDesc, func(t *testing.T) {
			testDep := provideTest(t)
			defer testDep()
			tc.mockFunc()

			got, err := svc.GetWalletBalance(context.Background(), tc.args.customerXID)
			assert.Equal(t, err != nil, tc.wantErr)
			assert.Equal(t, got, tc.wantResult)
		})
	}
}

func TestEnableWallet(t *testing.T) {
	type (
		args struct {
			customerXID string
		}
	)

	testCases := []struct {
		testID     int
		testDesc   string
		args       args
		mockFunc   func()
		wantErr    bool
		wantResult web.WalletResponse
	}{
		{
			testID:   1,
			testDesc: "Success",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "disabled",
				}, nil)
				mockRepository.EXPECT().UpdateWalletStatus(gomock.Any(), gomock.Any(), "enabled", gomock.Any()).Return(nil)
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "enabled",
				}, nil)
			},
			wantErr: false,
			wantResult: web.WalletResponse{
				ID:     "mock-id",
				Status: "enabled",
			},
		},
		{
			testID:   2,
			testDesc: "Failed - error get wallet",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{}, fmt.Errorf("error"))
			},
			wantErr:    true,
			wantResult: web.WalletResponse{},
		},
		{
			testID:   3,
			testDesc: "Failed - status enabled",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "enabled",
				}, nil)
			},
			wantErr:    true,
			wantResult: web.WalletResponse{},
		},
		{
			testID:   4,
			testDesc: "Failed - error UpdateWalletStatus",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "disabled",
				}, nil)
				mockRepository.EXPECT().UpdateWalletStatus(gomock.Any(), gomock.Any(), "enabled", gomock.Any()).Return(fmt.Errorf("error"))
			},
			wantErr:    true,
			wantResult: web.WalletResponse{},
		},
		{
			testID:   5,
			testDesc: "Failed - error get wallet",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "disabled",
				}, nil)
				mockRepository.EXPECT().UpdateWalletStatus(gomock.Any(), gomock.Any(), "enabled", gomock.Any()).Return(nil)
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{}, fmt.Errorf("error"))
			},
			wantErr:    true,
			wantResult: web.WalletResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDesc, func(t *testing.T) {
			testDep := provideTest(t)
			defer testDep()
			tc.mockFunc()

			got, err := svc.EnableWallet(context.Background(), tc.args.customerXID)
			assert.Equal(t, err != nil, tc.wantErr)
			assert.Equal(t, got, tc.wantResult)
		})
	}
}

func TestDisableWallet(t *testing.T) {
	type (
		args struct {
			customerXID string
		}
	)

	testCases := []struct {
		testID     int
		testDesc   string
		args       args
		mockFunc   func()
		wantErr    bool
		wantResult web.WalletResponse
	}{
		{
			testID:   1,
			testDesc: "Success",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().UpdateWalletStatus(gomock.Any(), gomock.Any(), "disabled", gomock.Any()).Return(nil)
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "disabled",
				}, nil)
			},
			wantErr: false,
			wantResult: web.WalletResponse{
				ID:     "mock-id",
				Status: "disabled",
			},
		},
		{
			testID:   2,
			testDesc: "Failed - error UpdateWalletStatus",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().UpdateWalletStatus(gomock.Any(), gomock.Any(), "disabled", gomock.Any()).Return(fmt.Errorf("error"))
			},
			wantErr:    true,
			wantResult: web.WalletResponse{},
		},
		{
			testID:   3,
			testDesc: "Failed - error get wallet",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().UpdateWalletStatus(gomock.Any(), gomock.Any(), "disabled", gomock.Any()).Return(nil)
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{}, fmt.Errorf("error"))
			},
			wantErr:    true,
			wantResult: web.WalletResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDesc, func(t *testing.T) {
			testDep := provideTest(t)
			defer testDep()
			tc.mockFunc()

			got, err := svc.DisableWallet(context.Background(), tc.args.customerXID)
			assert.Equal(t, err != nil, tc.wantErr)
			assert.Equal(t, got, tc.wantResult)
		})
	}
}

func TestWalletTransactions(t *testing.T) {
	type (
		args struct {
			customerXID string
		}
	)

	testCases := []struct {
		testID     int
		testDesc   string
		args       args
		mockFunc   func()
		wantErr    bool
		wantResult []web.TransactionResponse
	}{
		{
			testID:   1,
			testDesc: "Success",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "enabled",
				}, nil)
				mockRepository.EXPECT().GetWalletTransactions(gomock.Any(), "mock-id").Return([]domain.Transaction{
					{
						ID:              "mock-id-1",
						Amount:          20000,
						TransactionType: "deposit",
					},
					{
						ID:              "mock-id-2",
						Amount:          10000,
						TransactionType: "withdrawal",
					},
				}, nil)

			},
			wantErr: false,
			wantResult: []web.TransactionResponse{
				{
					ID:     "mock-id-1",
					Amount: 20000,
					Type:   "deposit",
				},
				{
					ID:     "mock-id-2",
					Amount: 10000,
					Type:   "withdrawal",
				},
			},
		},
		{
			testID:   2,
			testDesc: "Failed - error GetWaller",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{}, fmt.Errorf("error"))

			},
			wantErr:    true,
			wantResult: []web.TransactionResponse{},
		},
		{
			testID:   3,
			testDesc: "Failed - wallet disabled",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "disabled",
				}, nil)
			},
			wantErr:    true,
			wantResult: []web.TransactionResponse{},
		},
		{
			testID:   4,
			testDesc: "Failed - error GetWalletTransactions",
			args: args{
				customerXID: "1",
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "enabled",
				}, nil)
				mockRepository.EXPECT().GetWalletTransactions(gomock.Any(), "mock-id").Return([]domain.Transaction{}, fmt.Errorf("error"))
			},
			wantErr:    true,
			wantResult: []web.TransactionResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDesc, func(t *testing.T) {
			testDep := provideTest(t)
			defer testDep()
			tc.mockFunc()

			got, err := svc.GetWalletTransactions(context.Background(), tc.args.customerXID)
			assert.Equal(t, err != nil, tc.wantErr)
			assert.Equal(t, got, tc.wantResult)
		})
	}
}

func TestAddWalletBalance(t *testing.T) {
	type (
		args struct {
			customerXID string
			payload     web.TransactionRequest
		}
	)

	testCases := []struct {
		testID     int
		testDesc   string
		args       args
		mockFunc   func()
		wantErr    bool
		wantResult web.DepositResponse
	}{
		{
			testID:   1,
			testDesc: "Success",
			args: args{
				customerXID: "1",
				payload: web.TransactionRequest{
					Amount:      1000,
					ReferenceID: "mock-ref",
				},
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "enabled",
				}, nil)
				mockRepository.EXPECT().AddTransaction(gomock.Any(), gomock.Any()).Return(nil)
				mockRepository.EXPECT().UpdateWalletBalance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockRepository.EXPECT().UpdateTransactionStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
			wantResult: web.DepositResponse{
				Amount:      1000,
				ReferenceID: "mock-ref",
			},
		},
		{
			testID:   2,
			testDesc: "Failed - error validate",
			args: args{
				customerXID: "1",
				payload: web.TransactionRequest{
					Amount:      0,
					ReferenceID: "mock-ref",
				},
			},
			mockFunc: func() {

			},
			wantErr:    true,
			wantResult: web.DepositResponse{},
		},
		{
			testID:   3,
			testDesc: "Failed - error GetWallet",
			args: args{
				customerXID: "1",
				payload: web.TransactionRequest{
					Amount:      1000,
					ReferenceID: "mock-ref",
				},
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "enabled",
				}, fmt.Errorf("error"))
			},
			wantErr:    true,
			wantResult: web.DepositResponse{},
		},
		{
			testID:   4,
			testDesc: "Failed - wallet disable",
			args: args{
				customerXID: "1",
				payload: web.TransactionRequest{
					Amount:      1000,
					ReferenceID: "mock-ref",
				},
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "disabled",
				}, nil)
			},
			wantErr:    true,
			wantResult: web.DepositResponse{},
		},
		{
			testID:   5,
			testDesc: "Failed - error AddTransaction",
			args: args{
				customerXID: "1",
				payload: web.TransactionRequest{
					Amount:      1000,
					ReferenceID: "mock-ref",
				},
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "enabled",
				}, nil)
				mockRepository.EXPECT().AddTransaction(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
			wantErr:    true,
			wantResult: web.DepositResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDesc, func(t *testing.T) {
			testDep := provideTest(t)
			defer testDep()
			tc.mockFunc()

			got, err := svc.AddWalletBalance(context.Background(), tc.args.customerXID, tc.args.payload)
			assert.Equal(t, err != nil, tc.wantErr)
			assert.Equal(t, got.Amount, tc.wantResult.Amount)
			assert.Equal(t, got.ReferenceID, tc.wantResult.ReferenceID)
		})
	}
}

func TestDeductWalletBalance(t *testing.T) {
	type (
		args struct {
			customerXID string
			payload     web.TransactionRequest
		}
	)

	testCases := []struct {
		testID     int
		testDesc   string
		args       args
		mockFunc   func()
		wantErr    bool
		wantResult web.WithdrawalResponse
	}{
		{
			testID:   1,
			testDesc: "Success",
			args: args{
				customerXID: "1",
				payload: web.TransactionRequest{
					Amount:      1000,
					ReferenceID: "mock-ref",
				},
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:      "mock-id",
					Status:  "enabled",
					Balance: 1000000,
				}, nil)
				mockRepository.EXPECT().AddTransaction(gomock.Any(), gomock.Any()).Return(nil)
				mockRepository.EXPECT().UpdateWalletBalance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockRepository.EXPECT().UpdateTransactionStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
			wantResult: web.WithdrawalResponse{
				Amount:      1000,
				ReferenceID: "mock-ref",
			},
		},
		{
			testID:   2,
			testDesc: "Failed - error validate",
			args: args{
				customerXID: "1",
				payload: web.TransactionRequest{
					Amount:      0,
					ReferenceID: "mock-ref",
				},
			},
			mockFunc: func() {

			},
			wantErr:    true,
			wantResult: web.WithdrawalResponse{},
		},
		{
			testID:   3,
			testDesc: "Failed - error GetWallet",
			args: args{
				customerXID: "1",
				payload: web.TransactionRequest{
					Amount:      1000,
					ReferenceID: "mock-ref",
				},
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:     "mock-id",
					Status: "enabled",
				}, fmt.Errorf("error"))
			},
			wantErr:    true,
			wantResult: web.WithdrawalResponse{},
		},
		{
			testID:   4,
			testDesc: "Failed - wallet disable",
			args: args{
				customerXID: "1",
				payload: web.TransactionRequest{
					Amount:      1000,
					ReferenceID: "mock-ref",
				},
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:      "mock-id",
					Status:  "disabled",
					Balance: 1000000,
				}, nil)
			},
			wantErr:    true,
			wantResult: web.WithdrawalResponse{},
		},
		{
			testID:   4,
			testDesc: "Failed - insufficient balance",
			args: args{
				customerXID: "1",
				payload: web.TransactionRequest{
					Amount:      1000,
					ReferenceID: "mock-ref",
				},
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:      "mock-id",
					Status:  "enabled",
					Balance: 100,
				}, nil)
			},
			wantErr:    true,
			wantResult: web.WithdrawalResponse{},
		},
		{
			testID:   5,
			testDesc: "Failed - error AddTransaction",
			args: args{
				customerXID: "1",
				payload: web.TransactionRequest{
					Amount:      1000,
					ReferenceID: "mock-ref",
				},
			},
			mockFunc: func() {
				mockRepository.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Return(domain.Wallet{
					ID:      "mock-id",
					Status:  "enabled",
					Balance: 1000000,
				}, nil)
				mockRepository.EXPECT().AddTransaction(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
			wantErr:    true,
			wantResult: web.WithdrawalResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDesc, func(t *testing.T) {
			testDep := provideTest(t)
			defer testDep()
			tc.mockFunc()

			got, err := svc.DeductWalletBalance(context.Background(), tc.args.customerXID, tc.args.payload)
			assert.Equal(t, err != nil, tc.wantErr)
			assert.Equal(t, got.Amount, tc.wantResult.Amount)
			assert.Equal(t, got.ReferenceID, tc.wantResult.ReferenceID)
		})
	}
}
