package http

import (
	"github.com/gorilla/mux"
	v1 "github.com/rtzgod/EWallet/internal/delivery/http/v1"
	"github.com/rtzgod/EWallet/internal/domain/service"
)

type Service interface {
	GetWalletService() *service.WalletService
	GetTransactionService() *service.TransactionService
}

type Handlers struct {
	Transaction *v1.TransactionHandler
	Wallet      *v1.WalletHandler
}

func NewHandler(service Service) *Handlers {
	return &Handlers{
		Transaction: v1.NewTransactionHandler(service.GetTransactionService()),
		Wallet:      v1.NewWalletHandler(service.GetWalletService()),
	}
}

func (h *Handlers) InitRoutes(r *mux.Router) {
	h.Transaction.InitRoutes(r)
	h.Wallet.InitRoutes(r)
}
