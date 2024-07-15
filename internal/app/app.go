package app

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rtzgod/EWallet/internal/delivery/http/v1"
	"github.com/rtzgod/EWallet/internal/domain/service"
	"github.com/rtzgod/EWallet/internal/repository/psql"
	"log"
	"net/http"
	"os"
)

func Run() {
	r := mux.NewRouter()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("SERVER_PORT")

	walletStorage, transactionStorage := psql.NewWalletStorage(), psql.NewTransactionStorage()

	walletService := service.NewWalletService(walletStorage)
	transactionService := service.NewTransactionService(transactionStorage, walletService)

	wHandler := v1.NewWalletHandler(walletService)
	tHandler := v1.NewTransactionHandler(transactionService)

	wHandler.InitRoutes(r)
	tHandler.InitRoutes(r)

	done := make(chan bool)
	go http.ListenAndServe(":"+port, r)
	log.Println("Server started on port: " + port)
	<-done
}
