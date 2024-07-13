package service

import (
	"EWallet/internal/domain/entity"
	"EWallet/internal/repository/psql"
	"errors"
	"github.com/google/uuid"
	"sync"
)

var mu sync.Mutex
var db = psql.Connect()

func CreateWallet() *entity.Wallet {
	mu.Lock()
	defer mu.Unlock()
	id := generateID()
	wallet := &entity.Wallet{ID: id, Balance: 100.0}
	err := psql.AddWallet(db, id)
	if err != nil {
		panic(err)
	}
	return wallet
}

func GetWallet(id string) (*entity.Wallet, error) {
	mu.Lock()
	defer mu.Unlock()
	row, err := psql.GetWallets(db, id)
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
