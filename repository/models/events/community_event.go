package models

import "time"

type Community_Event struct {
	Id int
	Name string
	Percentage float64
	Code string
	Capacity int
	Date_Time time.Time
	Manager_Id int
	Event_Image *string
	Duration time.Duration
}