package service

import (
	"EWallet/internal/domain/entity"
	"EWallet/internal/repository/psql"
	"errors"
	"time"
)

func SendMoney(fromID, toID string, amount float64) error {
	senderWallet, senderExists := GetWallet(fromID)
	receiverWallet, receiverExists := GetWallet(toID)
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
	err := psql.UpdateBalance(db, fromID, amount)
	if err != nil {
		return err
	}
	err = psql.UpdateBalance(db, toID, amount)
	if err != nil {
		return err
	}
	transaction := entity.Transaction{
		Time:   time.Now(),
		From:   senderWallet.ID,
		To:     receiverWallet.ID,
		Amount: amount,
	}
	err = psql.AddTransaction(db, transaction)
	if err != nil {
		return err
	}
	return nil
}

func GetHistory(id string) ([]entity.Transaction, error) {
	_, err := GetWallet(id)
	if err != nil {
		return nil, errors.New("wallet not found")
	}
	mu.Lock()
	defer mu.Unlock()
	var transactions []entity.Transaction
	rows, err := psql.GetTransactions(db, id)
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
