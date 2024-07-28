package entity

import "time"

type Transaction struct {
	Time       time.Time `json:"time" db:"time"`
	SenderId   string    `json:"from" db:"sender_id"`
	ReceiverId string    `json:"to" db:"receiver_id" binding:"required"`
	Amount     float64   `json:"amount" db:"amount" binding:"required"`
}
