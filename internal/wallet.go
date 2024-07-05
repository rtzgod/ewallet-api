package internal

import (
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Wallet struct {
	ID      string
	Balance float64
}

type Transaction struct {
	Time   time.Time
	From   string
	To     string
	Amount float64
}

var wallets = make(map[string]*Wallet)
var transactions = make(map[string][]Transaction)
var mu sync.Mutex

func createWallet() map[string]*Wallet {
	mu.Lock()
	defer mu.Unlock()

	id := generateID()
	wallet := &Wallet{ID: id, Balance: 100.0}
	wallets[id] = wallet
	return wallets
}

func getWallet(id string) (*Wallet, error) {
	mu.Lock()
	defer mu.Unlock()
	wallet, exists := wallets[id]
	if !exists {
		return nil, errors.New("wallet not found")
	}

	return wallet, nil
}

func generateID() string {
	id := uuid.New()
	return id.String()
}
