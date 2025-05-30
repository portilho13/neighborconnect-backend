package controllers_models

import "time"

type ListingCreation struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Buy_Now_Price   string `json:"buy_now_price"`
	Start_Price     string `json:"start_price"`
	Expiration_Date string `json:"expiration_date"`
	Seller_Id       string `json:"seller_id"`
	Category_Id     string `json:"category_id"`
}

type Listing_Photos struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

type ListingInfo struct {
	Id              int               `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Buy_Now_Price   float64           `json:"buy_now_price"`
	Start_Price     float64           `json:"start_price"`
	Current_bid     BidInfo           `json:"current_bid"`
	Created_At      time.Time         `json:"created_at"`
	Expiration_Date time.Time         `json:"expiration_date"`
	Status          string            `json:"status"`
	Seller          SellerListingInfo `json:"seller"`
	Category        CategoryInfo      `json:"category"`
	Listing_Photos  []Listing_Photos  `json:"listing_photos"`
}
