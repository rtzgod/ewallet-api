package internal

import "time"

type Wallet struct {
	ID      string
	Balance float64
}

type Transaction struct {
	Time   time.Time
	From   string
	To     string
	Amount float64
}

var wallets = make(map[string]*Wallet)
var transactions = make(map[string][]Transaction)
