package v1

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"net/http"
)

type TransactionService interface {
	SendMoney(fromID, toID string, amount float64) error
	GetHistory(id string) ([]entity.Transaction, error)
}

type TransactionHandler struct {
	service TransactionService
}

func NewTransactionHandler(service TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) InitRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/wallet/{walletId}/send", h.SendMoney).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/history", h.GetHistory).Methods("GET")
}

func (h *TransactionHandler) SendMoney(w http.ResponseWriter, r *http.Request) {
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

func (h *TransactionHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	history, err := h.service.GetHistory(params["walletId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(history)
}
