package models

import "time"

type Rent struct {
	Id           *int
	Month        int
	Year         int
	Base_Amount  float64
	Reduction    float64
	Final_Amount float64
	Apartment_Id *int
	Status       string
	Due_day      int
	Pay_Day time.Time
}
