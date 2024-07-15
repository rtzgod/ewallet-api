package psql

import "database/sql"

func NewRepository() *Repository {
	db := Connect()
	return &Repository{db: db}
}

type Repository struct {
	db *sql.DB
}
