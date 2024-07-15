package service

import (
	"errors"
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"time"
)

func (s *Service) SendMoney(fromID, toID string, amount float64) error {
	senderWallet, senderExists := s.GetWallet(fromID)
	receiverWallet, receiverExists := s.GetWallet(toID)
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
	err := s.repo.UpdateBalance(fromID, amount)
	if err != nil {
		return err
	}
	err = s.repo.UpdateBalance(toID, amount)
	if err != nil {
		return err
	}
	transaction := entity.Transaction{
		Time:   time.Now(),
		From:   senderWallet.ID,
		To:     receiverWallet.ID,
		Amount: amount,
	}
	err = s.repo.AddTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetHistory(id string) ([]entity.Transaction, error) {
	_, err := s.GetWallet(id)
	if err != nil {
		return nil, errors.New("wallet not found")
	}
	mu.Lock()
	defer mu.Unlock()
	var transactions []entity.Transaction
	rows, err := s.repo.GetTransactions(id)
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
