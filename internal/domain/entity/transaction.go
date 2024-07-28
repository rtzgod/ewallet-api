package entity

import "time"

type Transaction struct {
	Time       time.Time `json:"-"`
	SenderId   string    `json:"-"`
	ReceiverId string    `json:"to" binding:"required"`
	Amount     float64   `json:"amount" binding:"required"`
}
