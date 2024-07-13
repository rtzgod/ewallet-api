package psql

import (
	"EWallet/internal/domain/entity"
	"database/sql"
	_ "github.com/lib/pq"
)

func AddTransaction(db *sql.DB, transaction entity.Transaction) error {
	_, err := db.Exec("insert into transactions (time, senderid, receiverid, amount) values ($1, $2, $3, $4)", transaction.Time, transaction.From, transaction.To, transaction.Amount)

	return err
}
func GetTransactions(db *sql.DB, id string) (*sql.Rows, error) {
	rows, err := db.Query("select * from transactions where senderid = $1 or receiverid = $1", id)

	return rows, err
}
