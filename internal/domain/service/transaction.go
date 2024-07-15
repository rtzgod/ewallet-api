package service

import (
	"database/sql"
	"errors"
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"time"
)

type TransactionStorage interface {
	AddTransaction(transaction entity.Transaction) error
	GetTransactions(id string) (*sql.Rows, error)
	UpdateBalance(id string, amount float64) error
}

type WalletProvider interface {
	GetWallet(id string) (*entity.Wallet, error)
}

type transactionService struct {
	storage        TransactionStorage
	walletProvider WalletProvider
}

func NewTransactionService(storage TransactionStorage, walletProvider WalletProvider) *transactionService {
	return &transactionService{storage: storage, walletProvider: walletProvider}
}

func (t *transactionService) SendMoney(fromID, toID string, amount float64) error {
	senderWallet, senderExists := t.walletProvider.GetWallet(fromID)
	receiverWallet, receiverExists := t.walletProvider.GetWallet(toID)
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
	err := t.storage.UpdateBalance(fromID, amount)
	if err != nil {
		return err
	}
	err = t.storage.UpdateBalance(toID, amount)
	if err != nil {
		return err
	}
	transaction := entity.Transaction{
		Time:   time.Now(),
		From:   senderWallet.ID,
		To:     receiverWallet.ID,
		Amount: amount,
	}
	err = t.storage.AddTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (t *transactionService) GetHistory(id string) ([]entity.Transaction, error) {
	_, err := t.walletProvider.GetWallet(id)
	if err != nil {
		return nil, errors.New("wallet not found")
	}
	mu.Lock()
	defer mu.Unlock()
	var transactions []entity.Transaction
	rows, err := t.storage.GetTransactions(id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		t := entity.Transaction{}
		err = rows.Scan(&t.Time, &t.From, &t.To, &t.Amount)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}
