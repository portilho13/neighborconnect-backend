package models

import "time"

type Listing struct {
	Id              *int
	Name            string
	Description     string
	Buy_Now_Price   int
	Start_Price     int
	Created_At      time.Time
	Expiration_Date time.Time
	Status          string
	Seller_Id       *int
	Category_Id     *int
}
