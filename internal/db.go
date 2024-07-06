package internal

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func Connect() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	userName := os.Getenv("USER_NAME")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DATABASE")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", userName, password, dbname)
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Exec("create table if not exists transactions (time timestamp, senderID varchar(100), receiverID varchar(100), amount numeric)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Exec("create table if not exists wallets (id varchar(100) primary key, balance numeric)")
	if err != nil {
		log.Fatal(err)
	}
	return database
}
