package controllers_models

import "time"

type AccountJson struct {
	Id            int     `json:"id"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency"`
	UsersID       *int    `json:"users_id"`
}

type AccountMovementJson struct {
	Id         int       `json:"id"`
	Amount     float64   `json:"amount"`
	Created_At time.Time `json:"created_at"`
	Account_Id *int      `json:"account_id"`
	Type       string    `json:"type"`
}
