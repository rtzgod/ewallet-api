package psql

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type WalletStorage struct {
	db *sql.DB
}

func NewWalletStorage(db *sql.DB) *WalletStorage {
	return &WalletStorage{db: db}
}

func (s *WalletStorage) AddWallet(id string) error {
	_, err := s.db.Exec("insert into wallets (id, balance) values ($1, 100)", id)

	return err
}

func (s *WalletStorage) GetWallets(id string) (*sql.Rows, error) {
	rows, err := s.db.Query("select * from wallets where id = $1", id)

	return rows, err
}
