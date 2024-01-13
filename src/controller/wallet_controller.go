package controller

import (
	"net/http"
)

type WalletController interface {
	InitializeWallet(writer http.ResponseWriter, request *http.Request)
	EnableWallet(writer http.ResponseWriter, request *http.Request)
	GetWalletBalance(writer http.ResponseWriter, request *http.Request)
	GetWalletTransactions(writer http.ResponseWriter, request *http.Request)
	AddMoneyToWallet(writer http.ResponseWriter, request *http.Request)
	WithdrawFromWallet(writer http.ResponseWriter, request *http.Request)
	DisableWallet(writer http.ResponseWriter, request *http.Request)
}
