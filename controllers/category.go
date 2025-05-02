package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func GetCategories(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	categories, err := repositoryControllers.GetAllCategories(dbPool)
	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	var categoriesJson []controllers_models.CategoryInfo

	for _, category := range categories {
		var url string
		if category.Url == nil {
			url = ""
		} else {
			url = *category.Url
		}
		categoryJson := controllers_models.CategoryInfo{
			Id:   *category.Id,
			Name: category.Name,
			Url:  url,
		}

		categoriesJson = append(categoriesJson, categoryJson)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(categoriesJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func CreateCategory(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var category controllers_models.CategoryCreation
	err := json.NewDecoder(r.Body).Decode(&category)

	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	categoryDb := models.Category{
		Name: category.Name,
		Url:  &category.Url,
	}

	err = repositoryControllers.CreateCategory(categoryDb, dbPool)
	if err != nil {
		http.Error(w, "Error Creating Category", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Category Created !"})

}
