package models

type User struct {
	Id int
	Name     string
	Email    string
	Password string
	Phone    string
	Apartment_id *int
	Profile_Picture *string
}