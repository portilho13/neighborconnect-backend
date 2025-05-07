package models

import "time"

type Transaction struct {
	Id               *int
	Final_Price      float64
	Transaction_Time time.Time
	Transaction_Type string
	Seller_Id        *int
	Buyer_Id         *int
	Listing_Id       *int
	Payment_Status   string
	Payment_Due_time time.Time
}
