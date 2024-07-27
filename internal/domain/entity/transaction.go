package entity

import "time"

type Transaction struct {
	Time       time.Time
	SenderId   string
	ReceiverId string  `json:"to" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
}
