package controllers_models

type AccountJson struct {
	Id            int     `json:"id"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency"`
	UsersID       *int    `json:"users_id"`
}
