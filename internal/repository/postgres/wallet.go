package postgres

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rtzgod/EWallet/internal/domain/entity"
)

type WalletPostgres struct {
	db *sqlx.DB
}

const baseAmount = 100.0

func NewWalletPostgres(db *sqlx.DB) *WalletPostgres {
	return &WalletPostgres{db: db}
}

func (r *WalletPostgres) Create(id string) (entity.Wallet, error) {
	var wallet entity.Wallet

	query := fmt.Sprintf("insert into %s (id, amount) values ($1, $2)", walletsTable)

	if _, err := r.db.Exec(query, id, baseAmount); err != nil {
		return wallet, errors.New("failed to insert into wallet: " + err.Error())
	}

	wallet = entity.Wallet{
		Id:     id,
		Amount: baseAmount,
	}

	return wallet, nil
}
