package controllers_models

type Withdraw struct {
	Type    string  `json:"type"`
	Amount  float64 `json:"amount"`
	User_id int     `json:"user_id"`
}
