package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rtzgod/ewallet-api/internal/domain/entity"
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
	_, err := r.db.Exec(query, time.Now().Format(time.RFC3339), senderId, receiverId, amount)
	return err
}

func (r *TransactionPostgres) GetAllById(id string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	query := fmt.Sprintf("select * from %s where sender_id = $1 or receiver_id = $2", transactionsTable)
	err := r.db.Select(&transactions, query, id, id)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
