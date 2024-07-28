package service

import (
	"github.com/google/uuid"
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"github.com/rtzgod/EWallet/internal/repository"
)

type WalletService struct {
	repo repository.Wallet
}

func NewWalletService(repo repository.Wallet) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) Create() (entity.Wallet, error) {
	id := generateId()
	wallet, err := s.repo.Create(id)
	if err != nil {
		return entity.Wallet{}, err
	}
	return wallet, nil
}

func (s *WalletService) GetById(id string) (entity.Wallet, error) {
	wallet, err := s.repo.GetById(id)
	if err != nil {
		return entity.Wallet{}, err
	}
	return wallet, nil
}

func generateId() string {
	id := uuid.New()
	return id.String()
}
