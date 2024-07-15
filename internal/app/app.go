package app

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rtzgod/EWallet/internal/delivery/httpv1"
	service2 "github.com/rtzgod/EWallet/internal/domain/service"
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

	repo := psql.NewRepository()
	service := service2.NewService(repo)
	h := httpv1.NewHandler(service)

	r.HandleFunc("/api/v1/wallet", h.CreateWallet).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/send", h.SendMoney).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/history", h.GetHistory).Methods("GET")
	r.HandleFunc("/api/v1/wallet/{walletId}", h.GetWallet).Methods("GET")

	done := make(chan bool)
	go http.ListenAndServe(":"+port, r)
	log.Println("Server started on port: " + port)
	<-done
}
