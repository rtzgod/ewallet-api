package psql

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rtzgod/EWallet/internal/domain/entity"
)

type TransactionStorage struct {
	db *sql.DB
}

func NewTransactionStorage(db *sql.DB) *TransactionStorage {
	return &TransactionStorage{db: db}
}

func (s *TransactionStorage) AddTransaction(transaction entity.Transaction) error {
	_, err := s.db.Exec("insert into transactions (time, senderid, receiverid, amount) values ($1, $2, $3, $4)", transaction.Time, transaction.From, transaction.To, transaction.Amount)

	return err
}
func (s *TransactionStorage) GetTransactions(id string) (*sql.Rows, error) {
	rows, err := s.db.Query("select * from transactions where senderid = $1 or receiverid = $1", id)

	return rows, err
}
func (s *TransactionStorage) UpdateBalance(id string, amount float64) error {
	_, err := s.db.Exec("update wallets set balance = balance - $1 where id = $2", amount, id)

	return err
}
