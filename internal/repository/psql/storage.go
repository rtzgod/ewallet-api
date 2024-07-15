package psql

import "database/sql"

type Storage struct {
	db *sql.DB
}

func NewStorage() *Storage {
	db := Connect()
	return &Storage{db: db}
}
