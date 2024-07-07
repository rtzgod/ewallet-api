package main

import (
	inter "EWallet/internal"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("SERVER_PORT")
	r.HandleFunc("/api/v1/wallet", inter.CreateWallet).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/send", inter.SendMoney).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/history", inter.GetHistory).Methods("GET")
	r.HandleFunc("/api/v1/wallet/{walletId}", inter.GetWallet).Methods("GET")
	fmt.Println("Connecting...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
