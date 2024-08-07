package postgres

import (
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

	query := fmt.Sprintf("insert into %s (id, balance) values ($1, $2)", walletsTable)

	if _, err := r.db.Exec(query, id, baseAmount); err != nil {
		return wallet, err
	}

	wallet = entity.Wallet{
		Id:      id,
		Balance: baseAmount,
	}

	return wallet, nil
}

func (r *WalletPostgres) GetById(id string) (entity.Wallet, error) {
	var wallet entity.Wallet
	query := fmt.Sprintf("select * from %s where id=$1", walletsTable)
	err := r.db.Get(&wallet, query, id)
	if err != nil {
		return wallet, err
	}
	return wallet, nil
}

func (r *WalletPostgres) Update(senderId, receiverId string, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	updateSenderBalanceQuery := fmt.Sprintf("update %s set balance = balance - $1 where id = $2", walletsTable)
	_, err = r.db.Exec(updateSenderBalanceQuery, amount, senderId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	updateReceiverBalanceQuery := fmt.Sprintf("update %s set balance = balance + $1 where id = $2", walletsTable)
	_, err = r.db.Exec(updateReceiverBalanceQuery, amount, receiverId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
