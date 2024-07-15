package service

import (
	"database/sql"
	"github.com/rtzgod/EWallet/internal/domain/entity"
)

type Storage interface {
	AddWallet(id string) error
	GetWallets(id string) (*sql.Rows, error)
	UpdateBalance(id string, amount float64) error
	AddTransaction(transaction entity.Transaction) error
	GetTransactions(id string) (*sql.Rows, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}
