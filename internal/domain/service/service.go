package service

import "github.com/rtzgod/EWallet/internal/repository/psql"

type Storage interface {
	GetWalletStorage() *psql.WalletStorage
	GetTransactionStorage() *psql.TransactionStorage
}

type Service struct {
	WalletService      *WalletService
	TransactionService *TransactionService
}

func NewService(storage Storage) *Service {
	return &Service{
		WalletService:      NewWalletService(storage.GetWalletStorage()),
		TransactionService: NewTransactionService(storage.GetTransactionStorage(), NewWalletService(storage.GetWalletStorage())),
	}
}

func (s *Service) GetWalletService() *WalletService {
	return s.WalletService
}

func (s *Service) GetTransactionService() *TransactionService {
	return s.TransactionService
}
