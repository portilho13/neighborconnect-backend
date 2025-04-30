package controllers_models

type ManagerDashboardInfo struct {
	Apartments []Apartment `json:"apartments"`
	Users []UserLogin `json:"users"`
}
