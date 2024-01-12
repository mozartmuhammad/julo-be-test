package app

import (
	"github.com/gorilla/mux"
	"github.com/mozartmuhammad/julo-be-test/controller"
	"github.com/mozartmuhammad/julo-be-test/middleware"
)

func NewRouter(walletController controller.WalletController) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/init", walletController.InitializeWallet).Methods("POST")
	router.HandleFunc("/api/v1/wallet", middleware.AuthorizeRequest(walletController.EnableWallet)).Methods("POST")
	router.HandleFunc("/api/v1/wallet", middleware.AuthorizeRequest(walletController.GetWalletBalance)).Methods("GET")
	router.HandleFunc("/api/v1/wallet", middleware.AuthorizeRequest(walletController.DisableWallet)).Methods("PATCH")
	router.HandleFunc("/api/v1/wallet/transactions", middleware.AuthorizeRequest(walletController.GetWalletTransactions)).Methods("GET")
	router.HandleFunc("/api/v1/wallet/deposits", middleware.AuthorizeRequest(walletController.AddMoneyToWallet)).Methods("POST")
	router.HandleFunc("/api/v1/wallet/withdrawals", middleware.AuthorizeRequest(walletController.WithdrawFromWallet)).Methods("POST")

	return router
}
