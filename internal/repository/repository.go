package repository

import (
	"database/sql"
	"github.com/rtzgod/EWallet/internal/repository/psql"
)

type Storage struct {
	WalletStorage      *psql.WalletStorage
	TransactionStorage *psql.TransactionStorage
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		WalletStorage:      psql.NewWalletStorage(db),
		TransactionStorage: psql.NewTransactionStorage(db),
	}
}

func (s *Storage) GetWalletStorage() *psql.WalletStorage {
	return s.WalletStorage
}

func (s *Storage) GetTransactionStorage() *psql.TransactionStorage {
	return s.TransactionStorage
}
