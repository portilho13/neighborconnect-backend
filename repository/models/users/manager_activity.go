package models

import "time"

type Manager_Activity struct {
	Id          *int
	Type        string
	Description string
	Created_At  time.Time
	Manager_Id  int
}
