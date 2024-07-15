package httpv1

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) InitRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/wallet", h.CreateWallet).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/send", h.SendMoney).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/history", h.GetHistory).Methods("GET")
	r.HandleFunc("/api/v1/wallet/{walletId}", h.GetWallet).Methods("GET")
}

func (h *Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	wallets := h.service.CreateWallet()
	_ = json.NewEncoder(w).Encode(wallets)
}

func (h *Handler) GetWallet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	wallet, err := h.service.GetWallet(params["walletId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	_ = json.NewEncoder(w).Encode(wallet)
}

func (h *Handler) SendMoney(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	senderID := params["walletId"]
	var req struct {
		ReceiverID string  `json:"to"`
		Amount     float64 `json:"amount"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	err := h.service.SendMoney(senderID, req.ReceiverID, req.Amount)
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

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	history, err := h.service.GetHistory(params["walletId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(history)
}
