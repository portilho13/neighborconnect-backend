package controllers_models

type CategoryCreation struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type CategoryInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}
