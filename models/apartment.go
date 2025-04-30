package controllers_models

type Apartment struct {
	Id         int `json:"id"`
	N_Bedrooms int `json:"n_bedrooms"`
	Floor      int `json:"floor"`
	Rent       int `json:"rent"`
	Manager_Id int `json:"manager_id"`
}
