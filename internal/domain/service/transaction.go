package service

import (
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"github.com/rtzgod/EWallet/internal/repository"
)

type TransactionService struct {
	repo repository.Transaction
}

func NewTransactionService(repo repository.Transaction) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Create(senderId, receiverId string, amount float64) error {
	return s.repo.Create(senderId, receiverId, amount)
}

func (s *TransactionService) GetAllById(id string) ([]entity.Transaction, error) {
	return s.repo.GetAllById(id)
}
