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
var db = Connect()

func createWallet() *Wallet {
	mu.Lock()
	defer mu.Unlock()
	id := generateID()
	wallet := &Wallet{ID: id, Balance: 100.0}
	_, err := db.Exec("insert into Wallets (id, balance) values ($1, 100)", id)
	if err != nil {
		panic(err)
	}
	return wallet
}

func getWallet(id string) (*Wallet, error) {
	mu.Lock()
	defer mu.Unlock()
	row, err := db.Query("select * from wallets where id = $1", id)
	defer row.Close()
	if err != nil {
		return nil, err
	}
	wallet := &Wallet{}
	row.Next()
	err = row.Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func sendMoney(fromID, toID string, amount float64) error {
	mu.Lock()
	defer mu.Unlock()
	senderWallet, fromExists := wallets[fromID]
	receiverWallet, toExists := wallets[toID]
	if !fromExists {
		return errors.New("sender wallet not found")
	}
	if !toExists {
		return errors.New("receiver wallet not found")
	}
	if senderWallet.Balance < amount {
		return errors.New("balance not enough")
	}
	if amount < 0 {
		return errors.New("amount is negative")
	}
	senderWallet.Balance -= amount
	receiverWallet.Balance += amount

	transaction := Transaction{
		Time:   time.Now(),
		From:   senderWallet.ID,
		To:     receiverWallet.ID,
		Amount: amount,
	}
	transactions[fromID] = append(transactions[fromID], transaction)
	transactions[toID] = append(transactions[toID], transaction)

	return nil
}

func getHistory(id string) ([]Transaction, error) {
	mu.Lock()
	defer mu.Unlock()
	_, exists := wallets[id]
	if !exists {
		return nil, errors.New("wallet not found")
	}
	return transactions[id], nil
}

func generateID() string {
	id := uuid.New()
	return id.String()
}
