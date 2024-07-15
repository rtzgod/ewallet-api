package service

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"sync"
)

type WalletStorage interface {
	AddWallet(id string) error
	GetWallets(id string) (*sql.Rows, error)
}

type walletService struct {
	storage WalletStorage
}

func NewWalletService(storage WalletStorage) *walletService {
	return &walletService{storage: storage}
}

var mu sync.Mutex

func (w *walletService) CreateWallet() *entity.Wallet {
	mu.Lock()
	defer mu.Unlock()
	id := generateID()
	wallet := &entity.Wallet{ID: id, Balance: 100.0}
	err := w.storage.AddWallet(id)
	if err != nil {
		panic(err)
	}
	return wallet
}

func (w *walletService) GetWallet(id string) (*entity.Wallet, error) {
	mu.Lock()
	defer mu.Unlock()
	row, err := w.storage.GetWallets(id)
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
