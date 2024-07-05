package internal

import (
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

func CreateWallet() *Wallet {
	mu.Lock()
	defer mu.Unlock()

	id := generateID()
	wallet := &Wallet{ID: id, Balance: 100.0}
	wallets[id] = wallet
	return wallet
}

func generateID() string {
	id := uuid.New()
	return id.String()
}
