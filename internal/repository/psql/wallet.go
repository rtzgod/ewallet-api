package psql

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type walletStorage struct {
	db *sql.DB
}

func NewWalletStorage() *walletStorage {
	db := Connect()
	return &walletStorage{db: db}
}

func (s *walletStorage) AddWallet(id string) error {
	_, err := s.db.Exec("insert into wallets (id, balance) values ($1, 100)", id)

	return err
}

func (s *walletStorage) GetWallets(id string) (*sql.Rows, error) {
	rows, err := s.db.Query("select * from wallets where id = $1", id)

	return rows, err
}
