package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"github.com/rtzgod/EWallet/internal/repository/postgres"
)

type Wallet interface {
	Create(id string) (entity.Wallet, error)
	GetById(id string) (entity.Wallet, error)
}

type Transaction interface {
}

type Repository struct {
	Wallet
	Transaction
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Wallet: postgres.NewWalletPostgres(db),
	}
}
