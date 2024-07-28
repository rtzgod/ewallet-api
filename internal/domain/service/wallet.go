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
	return s.repo.Create(id)
}

func (s *WalletService) GetById(id string) (entity.Wallet, error) {
	return s.repo.GetById(id)
}

func (s *WalletService) Update(senderId, receiverId string, amount float64) error {
	return s.repo.Update(senderId, receiverId, amount)
}

func generateId() string {
	id := uuid.New()
	return id.String()
}
