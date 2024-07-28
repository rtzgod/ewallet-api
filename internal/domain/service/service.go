package service

import (
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"github.com/rtzgod/EWallet/internal/repository"
)

type Wallet interface {
	Create() (entity.Wallet, error)
	GetById(id string) (entity.Wallet, error)
}

type Transaction interface {
}

type Service struct {
	Wallet
	Transaction
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Wallet: NewWalletService(repos.Wallet),
	}
}
