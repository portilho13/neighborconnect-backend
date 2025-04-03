package models

import "time"


type Bid struct {
	Id *int
	Bid_Ammount int
	Bid_Time time.Time
	Users_Id *int
	Listing_Id *int
}