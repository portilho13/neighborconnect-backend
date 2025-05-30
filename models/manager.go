package controllers_models

import "time"

type ManagerCreationJson struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type ManagerInfoJson struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ManagerActivity struct {
	Id          int       `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Created_At  time.Time `json:"created_at"`
}

type ManagerTransaction struct {
	Id          int       `json:"id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Users_Id    *int      `json:"users_id"`
}
