package models

import "time"

type Account_Movement struct {
	Id int
	Ammount float64
	Created_at time.Time
	Account_id *int
	Type string
}