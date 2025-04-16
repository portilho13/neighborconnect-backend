package controllers_models

type BidJson struct {
	Bid_Ammount int  `json:"bid_ammount"`
	User_Id     *int `json:"users_id"`
	Listing_Id  *int `json:"listing_id"`
}
