package psql

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rtzgod/EWallet/internal/domain/entity"
)

func (r *Repository) AddTransaction(transaction entity.Transaction) error {
	_, err := r.db.Exec("insert into transactions (time, senderid, receiverid, amount) values ($1, $2, $3, $4)", transaction.Time, transaction.From, transaction.To, transaction.Amount)

	return err
}
func (r *Repository) GetTransactions(id string) (*sql.Rows, error) {
	rows, err := r.db.Query("select * from transactions where senderid = $1 or receiverid = $1", id)

	return rows, err
}
