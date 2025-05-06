package controllers_models

import "time"

type BidCreation struct {
	Bid_Ammount float64 `json:"bid_ammount"`
	User_Id     *int    `json:"users_id"`
	Listing_Id  *int    `json:"listing_id"`
}

type BidInfo struct {
	Id          *int       `json:"id"`
	Bid_Ammount float64    `json:"bid_ammount"`
	Bid_Time    *time.Time `json:"bid_time"`
	User_Id     *int       `json:"users_id"`
	Listing_Id  int        `json:"listing_id"`
}
