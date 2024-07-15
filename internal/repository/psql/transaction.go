package psql

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rtzgod/EWallet/internal/domain/entity"
)

type transactionStorage struct {
	db *sql.DB
}

func NewTransactionStorage() *transactionStorage {
	db := Connect()
	return &transactionStorage{db: db}
}

func (s *transactionStorage) AddTransaction(transaction entity.Transaction) error {
	_, err := s.db.Exec("insert into transactions (time, senderid, receiverid, amount) values ($1, $2, $3, $4)", transaction.Time, transaction.From, transaction.To, transaction.Amount)

	return err
}
func (s *transactionStorage) GetTransactions(id string) (*sql.Rows, error) {
	rows, err := s.db.Query("select * from transactions where senderid = $1 or receiverid = $1", id)

	return rows, err
}
func (s *transactionStorage) UpdateBalance(id string, amount float64) error {
	_, err := s.db.Exec("update wallets set balance = balance - $1 where id = $2", amount, id)

	return err
}
