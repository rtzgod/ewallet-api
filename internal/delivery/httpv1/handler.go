package httpv1

import "github.com/rtzgod/EWallet/internal/domain/entity"

type Service interface {
	CreateWallet() *entity.Wallet
	GetWallet(id string) (*entity.Wallet, error)
	SendMoney(fromID, toID string, amount float64) error
	GetHistory(id string) ([]entity.Transaction, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}
