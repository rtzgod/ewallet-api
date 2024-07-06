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

var mu sync.Mutex
var db = Connect()

func createWallet() *Wallet {
	mu.Lock()
	defer mu.Unlock()
	id := generateID()
	wallet := &Wallet{ID: id, Balance: 100.0}
	_, err := db.Exec("insert into wallets (id, balance) values ($1, 100)", id)
	if err != nil {
		panic(err)
	}
	return wallet
}

func getWallet(id string) (*Wallet, error) {
	mu.Lock()
	defer mu.Unlock()
	row, err := db.Query("select * from wallets where id = $1", id)
	if err != nil {
		return nil, err
	}
	wallet := &Wallet{}
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

func sendMoney(fromID, toID string, amount float64) error {
	senderWallet, senderExists := getWallet(fromID)
	receiverWallet, receiverExists := getWallet(toID)
	mu.Lock()
	defer mu.Unlock()
	if senderExists != nil {
		return errors.New("sender wallet not found")
	}
	if receiverExists != nil {
		return errors.New("receiver wallet not found")
	}
	if senderWallet.Balance < amount {
		return errors.New("balance not enough")
	}
	if amount < 0 {
		return errors.New("amount is negative")
	}
	_, err := db.Exec("update wallets set balance = balance - $1 where id = $2", amount, fromID)
	if err != nil {
		return err
	}
	_, err = db.Exec("update wallets set balance = balance + $1 where id = $2", amount, toID)
	if err != nil {
		return err
	}
	transaction := Transaction{
		Time:   time.Now(),
		From:   senderWallet.ID,
		To:     receiverWallet.ID,
		Amount: amount,
	}
	_, err = db.Exec("insert into transactions (time, senderid, receiverid, amount) values ($1, $2, $3, $4)", transaction.Time, transaction.From, transaction.To, transaction.Amount)
	if err != nil {
		return err
	}
	return nil
}

func getHistory(id string) ([]Transaction, error) {
	_, err := getWallet(id)
	if err != nil {
		return nil, errors.New("wallet not found")
	}
	mu.Lock()
	defer mu.Unlock()
	var transactions []Transaction
	rows, err := db.Query("select * from transactions where senderid = $1 or receiverid = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		t := Transaction{}
		err = rows.Scan(&t.Time, &t.From, &t.To, &t.Amount)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func generateID() string {
	id := uuid.New()
	return id.String()
}
