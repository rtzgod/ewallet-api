package internal

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateWallet(w http.ResponseWriter, r *http.Request) {
	wallets := createWallet()
	json.NewEncoder(w).Encode(wallets)
}

func GetWallet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	wallet, err := getWallet(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	json.NewEncoder(w).Encode(wallet)
}
