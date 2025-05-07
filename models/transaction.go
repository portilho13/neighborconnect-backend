package controllers_models

type PayJson struct {
	Type           string `json:"type"`
	Transaction_Id int    `json:"transaction_id"`
	User_Id        int    `json:"user_id"`
}
