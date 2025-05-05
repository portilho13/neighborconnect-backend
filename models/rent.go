package controllers_models

import "time"

type Rent struct {
	Id           *int       `json:"id"`
	Month        int        `json:"month"`
	Year         int        `json:"year"`
	Base_Amount  float64    `json:"base_amount"`
	Reduction    float64    `json:"reduction"`
	Final_Amount float64    `json:"final_amount"`
	Apartment_Id *int       `json:"apartment_id"`
	Status       string     `json:"status"`
	Due_Day      int        `json:"due_day"`
	Pay_Day      *time.Time `json:"pay_day"`
}
