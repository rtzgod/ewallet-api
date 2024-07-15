package psql

import "database/sql"

func NewStorage() *Storage {
	db := Connect()
	return &Storage{db: db}
}

type Storage struct {
	db *sql.DB
}
