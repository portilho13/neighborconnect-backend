package controllers_models

import "time"

type PayJson struct {
	Type           string `json:"type"`
	Transaction_Id int    `json:"transaction_id"`
	User_Id        int    `json:"user_id"`
}

type TransactionJson struct {
	Id               *int      `json:"id"`
	Final_Price      float64   `json:"final_price"`
	Transaction_Time time.Time `json:"transaction_time"`
	Transaction_Type string    `json:"transaction_type"`
	Seller_Id        *int      `json:"seller_id"`
	Buyer_Id         *int      `json:"buyer_id"`
	Listing_Id       *int      `json:"listing_id"`
	Payment_Status   string    `json:"payment_status"`
	Payment_Due_time time.Time `json:"payment_due_time"`
}
