package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/wallet", createWallet).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/send", sendMoney).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/history", getHistory).Methods("GET")
	r.HandleFunc("/api/v1/wallet/{walletId}", getWallet).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
