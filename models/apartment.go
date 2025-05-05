package controllers_models

type ApartmentCreation struct {
	N_Bedrooms int     `json:"n_bedrooms"`
	Floor      int     `json:"floor"`
	Rent       float64 `json:"rent"`
	Manager_Id int     `json:"manager_id"`
}

type Apartment struct {
	Id         int    `json:"id"`
	N_Bedrooms int    `json:"n_bedrooms"`
	Floor      int    `json:"floor"`
	Rent       int    `json:"rent"`
	Manager_Id int    `json:"manager_id"`
	Status     string `json:"status"`
	Last_Rent  *Rent  `json:"last_rent"`
}
