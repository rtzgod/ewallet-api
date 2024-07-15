package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"sync"
)

var mu sync.Mutex

func (s *Service) CreateWallet() *entity.Wallet {
	mu.Lock()
	defer mu.Unlock()
	id := generateID()
	wallet := &entity.Wallet{ID: id, Balance: 100.0}
	err := s.repo.AddWallet(id)
	if err != nil {
		panic(err)
	}
	return wallet
}

func (s *Service) GetWallet(id string) (*entity.Wallet, error) {
	mu.Lock()
	defer mu.Unlock()
	row, err := s.repo.GetWallets(id)
	if err != nil {
		return nil, err
	}
	wallet := &entity.Wallet{}
	exists := row.Next()
	if !exists {
		return nil, errors.New("wallet not found")
	}
	err = row.Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func generateID() string {
	id := uuid.New()
	return id.String()
}
