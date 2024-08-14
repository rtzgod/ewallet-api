package service

import (
	"github.com/rtzgod/ewallet-api/internal/domain/entity"
	"github.com/rtzgod/ewallet-api/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Wallet interface {
	Create() (entity.Wallet, error)
	GetById(id string) (entity.Wallet, error)
	Update(senderId, receiverId string, amount float64) error
}

type Transaction interface {
	Create(senderId, receiverId string, amount float64) error
	GetAllById(id string) ([]entity.Transaction, error)
}

type Service struct {
	Wallet
	Transaction
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Wallet:      NewWalletService(repos.Wallet),
		Transaction: NewTransactionService(repos.Transaction),
	}
}
