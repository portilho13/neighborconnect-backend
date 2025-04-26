package controllers_models

type UserJson struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Phone       string `json:"phone"`
	ApartmentID int    `json:"apartment_id"`
}

type UserLogin struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	ApartmentID int    `json:"apartment_id"`
	Avatar      string `json:"avatar"`
}
