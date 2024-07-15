package v1

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"net/http"
)

type WalletService interface {
	CreateWallet() *entity.Wallet
	GetWallet(id string) (*entity.Wallet, error)
}

type WalletHandler struct {
	service WalletService
}

func NewWalletHandler(service WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) InitRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/wallet", h.CreateWallet).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}", h.GetWallet).Methods("GET")
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	wallets := h.service.CreateWallet()
	_ = json.NewEncoder(w).Encode(wallets)
}

func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	wallet, err := h.service.GetWallet(params["walletId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	_ = json.NewEncoder(w).Encode(wallet)
}
