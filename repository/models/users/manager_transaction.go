package models

import "time"

type Manager_Transaction struct {
	Id int
	Type string
	Amount float64
	Date time.Time
	Description string
	Users_Id *int
	Manager_Id int
}