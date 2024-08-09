package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	walletsTable      = "wallets"
	transactionsTable = "transactions"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgres(cfg Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping postgres")
	}
	if err := startMigration(db.DB); err != nil {
		return nil, err
	}
	return db, nil
}

func startMigration(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to create database driver")
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"postgres", driver)
	if err != nil {
		return errors.Wrap(err, "failed to create migration instance")
	}
	version, dirty, err := m.Version()
	if err != nil {
		return errors.Wrap(err, "failed to get migration version")
	}

	log.Printf("Current migration version: %d, dirty: %t", version, dirty)

	if version == 0 {
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return errors.Wrap(err, "failed to run migration")
		}
	}
	return nil
}
