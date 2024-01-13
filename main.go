package main

import (
	"net/http"
	"time"

	"github.com/mozartmuhammad/julo-be-test/src/app"
	"github.com/mozartmuhammad/julo-be-test/src/controller"
	"github.com/mozartmuhammad/julo-be-test/src/repository"
	"github.com/mozartmuhammad/julo-be-test/src/service"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := app.NewDB()
	validate := validator.New()
	walletRepository := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(walletRepository, validate, 5*time.Second)
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
