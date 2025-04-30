package controllers_models

import "time"

type ListingCreation struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Buy_Now_Price   int       `json:"buy_now_price"`
	Start_Price     int       `json:"start_price"`
	Expiration_Time time.Time `json:"expiration_time"`
	Seller_Id       int       `json:"seller_id"`
}

type ListingInfo struct {
	Id int `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Buy_Now_Price   int       `json:"buy_now_price"`
	Start_Price     int       `json:"start_price"`
	Current_bid 	int 	  `json:"current_bid"`
	Created_At      time.Time `json:"created_at"`
	Expiration_Time time.Time `json:"expiration_time"`
	Status          string    `json:"status"`
	Seller_Id       *int      `json:"seller_id"` // Remove this * in prod
	Category_Id *int `json:"category_id"`

}
