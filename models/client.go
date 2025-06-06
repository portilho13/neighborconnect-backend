package controllers_models

type UserCreation struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	Phone           string `json:"phone"`
	ApartmentID     int    `json:"apartment_id"`
	Profile_Picture string `json:"profile_picture"`
}

type UserLogin struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	ApartmentID int    `json:"apartment_id"`
	Avatar      string `json:"avatar"`
}

type SellerListingInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type NeighborInfo struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	ApartmentID     int    `json:"apartment_id"`
	Profile_Picture string `json:"profile_picture"`
}
