package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type TransactionPostgres struct {
	db *sqlx.DB
}

func NewTransactionPostgres(db *sqlx.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

func (r *TransactionPostgres) Create(senderId, receiverId string, amount float64) error {
	query := fmt.Sprintf("insert into %s (time, sender_id, receiver_id, amount) values ($1, $2, $3, $4)", transactionsTable)
	_, err := r.db.Exec(query, time.Now(), senderId, receiverId, amount)
	if err != nil {
		return err
	}
	return nil
}
