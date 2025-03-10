package models

import "time"

type Withdraw struct {
	Id int
	Ammount int
	Created_at time.Time
}