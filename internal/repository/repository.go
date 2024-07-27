package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rtzgod/EWallet/internal/domain/entity"
)

type Wallet interface {
	Create(Id string) entity.Wallet
}

type Transaction interface {
}

type Repository struct {
	Wallet
	Transaction
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
