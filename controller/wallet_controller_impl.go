package controller

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mozartmuhammad/julo-be-test/helper"
	"github.com/mozartmuhammad/julo-be-test/model/web"
	"github.com/mozartmuhammad/julo-be-test/service"
)

type WalletControllerImpl struct {
	WalletService service.WalletServiceItf
}

func NewWalletController(walletService service.WalletServiceItf) WalletController {
	return &WalletControllerImpl{
		WalletService: walletService,
	}
}

func (c *WalletControllerImpl) InitializeWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerXID := r.FormValue("customer_xid")
	walletCreateRequest := web.WalletCreateRequest{
		CustomerXID: customerXID,
	}

	err := c.WalletService.InitializeWallet(ctx, walletCreateRequest)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, _ := c.GenerateToken(customerXID)
	helper.WriteSuccess(w, map[string]interface{}{
		"token": token,
	})
}

func (c *WalletControllerImpl) EnableWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerXID := helper.GetCustomerXID(ctx)
	result, err := c.WalletService.EnableWallet(ctx, customerXID)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteSuccess(w, map[string]interface{}{
		"wallet": result,
	})
}

func (c *WalletControllerImpl) DisableWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerXID := helper.GetCustomerXID(ctx)

	isDisabledStr := r.FormValue("is_disabled")
	_, _ = strconv.ParseBool(isDisabledStr)

	result, err := c.WalletService.DisableWallet(ctx, customerXID)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteSuccess(w, map[string]interface{}{
		"wallet": result,
	})
}

func (c *WalletControllerImpl) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerXID := helper.GetCustomerXID(ctx)
	result, err := c.WalletService.GetWalletBalance(ctx, customerXID)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteSuccess(w, map[string]interface{}{
		"wallet": result,
	})
}

func (c *WalletControllerImpl) GetWalletTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerXID := helper.GetCustomerXID(ctx)

	result, err := c.WalletService.GetWalletTransactions(ctx, customerXID)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteSuccess(w, map[string]interface{}{
		"transactions": result,
	})
}

func (c *WalletControllerImpl) AddMoneyToWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerXID := helper.GetCustomerXID(ctx)

	referenceID := r.FormValue("reference_id")
	amountStr := r.FormValue("amount")
	amount, _ := strconv.Atoi(amountStr)

	result, err := c.WalletService.AddWalletBalance(ctx, customerXID, web.TransactionRequest{
		Amount:      amount,
		ReferenceID: referenceID,
	})
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteSuccess(w, map[string]interface{}{
		"deposit": result,
	})
}

func (c *WalletControllerImpl) WithdrawFromWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerXID := helper.GetCustomerXID(ctx)

	referenceID := r.FormValue("reference_id")
	amountStr := r.FormValue("amount")
	amount, _ := strconv.Atoi(amountStr)

	result, err := c.WalletService.DeductWalletBalance(ctx, customerXID, web.TransactionRequest{
		Amount:      amount,
		ReferenceID: referenceID,
	})
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteSuccess(w, map[string]interface{}{
		"withdrawal": result,
	})
}

func (c *WalletControllerImpl) GenerateToken(customerXID string) (token string, err error) {
	secret := os.Getenv("SECRET")
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_xid": customerXID,
		"exp":          time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err = claims.SignedString([]byte(secret))

	return
}
