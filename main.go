package main

import (
	"net/http"

	"github.com/mozartmuhammad/julo-be-test/app"
	"github.com/mozartmuhammad/julo-be-test/controller"
	"github.com/mozartmuhammad/julo-be-test/repository"
	"github.com/mozartmuhammad/julo-be-test/service"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := app.NewDB()
	validate := validator.New()
	walletRepository := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(walletRepository, db, validate)
	walletController := controller.NewWalletController(walletService)

	router := app.NewRouter(walletController)
	server := http.Server{
		Addr:    ":1323",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
