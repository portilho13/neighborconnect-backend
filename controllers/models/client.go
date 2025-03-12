package controllers_models

type UserJson struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Phone string `json:"string"`
	ApartmentID int `json:"apartment_id"`
}

