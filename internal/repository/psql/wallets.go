package psql

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func (r *Repository) AddWallet(id string) error {
	_, err := r.db.Exec("insert into wallets (id, balance) values ($1, 100)", id)

	return err
}

func (r *Repository) GetWallets(id string) (*sql.Rows, error) {
	rows, err := r.db.Query("select * from wallets where id = $1", id)

	return rows, err
}

func (r *Repository) UpdateBalance(id string, amount float64) error {
	_, err := r.db.Exec("update wallets set balance = balance - $1 where id = $2", amount, id)

	return err
}
