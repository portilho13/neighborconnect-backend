package controllers_models

import "time"

type EventCreation struct {
	Name       string        `json:"name"`
	Percentage float64       `json:"percentage"`
	Capacity   int           `json:"capacity"`
	Date_time  time.Time     `json:"date_time"`
	Manager_Id int           `json:"manager_id"`
	Duration   time.Duration `json:"duration"`
	Local      string        `json:"local"`
}

type EventInfo struct {
	Id                int           `json:"id"`
	Name              string        `json:"name"`
	Percentage        float64       `json:"percentage"`
	Capacity          int           `json:"capacity"`
	Date_time         time.Time     `json:"date_time"`
	Manager_Id        int           `json:"manager_id"`
	Event_Image       string        `json:"event_image"`
	Duration          time.Duration `json:"duration"`
	Local             string        `json:"local"`
	Current_Ocupation int           `json:"current_ocupation"`
}

type JoinEvent struct {
	Community_Event_Id int `json:"community_event_id"`
	User_Id            int `json:"user_id"`
}

type ConcludeEventJson struct {
	Event_Id          int   `json:"community_event_id"`
	Awarded_Users_Ids []int `json:"awarded_users_ids"`
}
