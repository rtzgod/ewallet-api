package main

import (
	inter "EWallet/internal"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/wallet", inter.CreateWallet).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/send", inter.SendMoney).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/history", inter.GetHistory).Methods("GET")
	r.HandleFunc("/api/v1/wallet/{walletId}", inter.GetWallet).Methods("GET")
	fmt.Println("server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
