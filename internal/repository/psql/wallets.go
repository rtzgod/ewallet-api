package psql

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func AddWallet(db *sql.DB, id string) error {
	_, err := db.Exec("insert into wallets (id, balance) values ($1, 100)", id)

	return err
}

func GetWallets(db *sql.DB, id string) (*sql.Rows, error) {
	rows, err := db.Query("select * from wallets where id = $1", id)

	return rows, err
}

func UpdateBalance(db *sql.DB, id string, amount float64) error {
	_, err := db.Exec("update wallets set balance = balance - $1 where id = $2", amount, id)

	return err
}
