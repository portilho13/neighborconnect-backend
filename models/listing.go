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
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Buy_Now_Price   int       `json:"buy_now_price"`
	Start_Price     int       `json:"start_price"`
	Created_At      time.Time `json:"created_at"`
	Expiration_Time time.Time `json:"expiration_time"`
	Status          string    `json:"status"`
	Seller_Id       *int      `json:"seller_id"` // Remove this in prod
}
