package models

import "time"

type Transaction struct {
	Id int
	Final_Price int
	Transaction_Time time.Time
	Transaction_Type string
	Seller_Id int
	Buyer_Id int
	Listning_Id *int
}