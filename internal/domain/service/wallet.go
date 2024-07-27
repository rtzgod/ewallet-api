package service

import (
	"github.com/google/uuid"
	"github.com/rtzgod/EWallet/internal/repository"
)

type WalletService struct {
	repo repository.Wallet
}

func NewWalletService(repo repository.Wallet) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) Create() {
	id := generateId()
	s.repo.Create(id)
}

func generateId() string {
	id := uuid.New()
	return id.String()
}
