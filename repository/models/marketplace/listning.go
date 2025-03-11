package models

import "time"

type Listning struct {
	Id int
	Start_Price int
	Buy_Now_Price int
	Expiration_Time time.Time
	Created_At time.Time
	Status string
	Seller_Id int
	Item_Id int
}