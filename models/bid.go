package controllers_models

type BidCreation struct {
	Bid_Ammount int  `json:"bid_ammount"`
	User_Id     *int `json:"users_id"`
	Listing_Id  *int `json:"listing_id"`
}

type BidInfo struct {
	Id          *int `json:"id"`
	Bid_Ammount int  `json:"bid_ammount"`
	User_Id     *int `json:"users_id"`
	Listing_Id  int  `json:"listing_id"`
}
