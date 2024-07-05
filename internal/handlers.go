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
	wallet, err := getWallet(params["walletId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	_ = json.NewEncoder(w).Encode(wallet)
}

func SendMoney(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	senderID := params["walletId"]
	var req struct {
		ReceiverID string  `json:"to"`
		Amount     float64 `json:"amount"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	err := sendMoney(senderID, req.ReceiverID, req.Amount)
	if err != nil {
		switch err.Error() {
		case "404":
			http.Error(w, err.Error(), http.StatusNotFound)
		case "400":
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
