package entity

import "time"

type Transaction struct {
	Time   time.Time
	From   string
	To     string
	Amount float64
}
