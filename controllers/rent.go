package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/portilho13/neighborconnect-backend/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetRents(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	apartment_id_str := r.URL.Query().Get("apartment_id")

	apartmend_id, err := strconv.Atoi(apartment_id_str)
	if err != nil {
		http.Error(w, "Invalid Apartment ID", http.StatusBadRequest)
		return
	}

	rents, err := repositoryControllers.GetRentByApartmentId(apartmend_id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching Rents", http.StatusInternalServerError)
		return
	}

	var rentsJson []controllers_models.Rent

	for _, rent := range rents {

		rentJson := controllers_models.Rent{
			Id:           rent.Id,
			Month:        rent.Month,
			Year:         rent.Year,
			Base_Amount:  rent.Base_Amount,
			Reduction:    rent.Reduction,
			Final_Amount: rent.Final_Amount,
			Apartment_Id: rent.Apartment_Id,
			Status:       rent.Status,
			Due_Day:      rent.Due_day,
			Pay_Day:      rent.Pay_Day,
		}

		rentsJson = append(rentsJson, rentJson)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rentsJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func PayRent(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var rentPayJson controllers_models.RentPay

	err := json.NewDecoder(r.Body).Decode(&rentPayJson)
	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	rent, err := repositoryControllers.GetRentById(rentPayJson.Rent_Id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching Rents", http.StatusInternalServerError)
		return
	}

	if rent.Status != "unpaid" {
		http.Error(w, "Rent Already Paid", http.StatusInternalServerError)
		return
	}

	user, err := repositoryControllers.GetUsersById(rentPayJson.User_Id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching User", http.StatusInternalServerError)
		return
	}

	switch rentPayJson.Payment_Type {
	case "wallet":
		if !utils.ValidateWalletBalance(user.Id, rent.Final_Amount, dbPool) {
			http.Error(w, "Insuficient Balance", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Unsuported Method", http.StatusInternalServerError)
		return

	}

	user_account, err := repositoryControllers.GetAccountByUserId(user.Id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching Seller Id Account", http.StatusInternalServerError)
		return
	}

	new_balance := user_account.Balance - rent.Final_Amount

	err = repositoryControllers.UpdateAccountBalance(user.Id, new_balance, dbPool)
	if err != nil {
		http.Error(w, "Error Updating Account Balance", http.StatusInternalServerError)
		return
	}

	err = repositoryControllers.UpdateRentStatus("paid", *rent.Id, dbPool)
	if err != nil {
		http.Error(w, "Error Updating Rent Status", http.StatusInternalServerError)
		return
	}

	err = repositoryControllers.UpdateRentPayday(time.Now(), *rent.Id, dbPool)
	if err != nil {
		http.Error(w, "Error Updating Rent Pay Day", http.StatusInternalServerError)
		return
	}

	manager_id, err := utils.GetManagerIdByUserId(user.Id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching Manager Id", http.StatusInternalServerError)
		return
	}
	description_str_activity := fmt.Sprintf("User %d payed their rent for apartment %d", user.Id, *user.Apartment_id)
	manager_activity := models.Manager_Activity{
		Type:        "rent",
		Description: description_str_activity,
		Created_At:  time.Now().UTC(),
		Manager_Id:  *manager_id,
	}

	err = repositoryControllers.CreateManagerActivity(manager_activity, dbPool)
	if err != nil {
		http.Error(w, "Error Creating Manager Activity", http.StatusInternalServerError)
		return
	}

	description_str_transaction := fmt.Sprintf("User %d payed their rent for apartment %d", user.Id, *user.Apartment_id)
	manager_transaction := models.Manager_Transaction{
		Type:        "rent",
		Amount:      rent.Final_Amount,
		Date:        time.Now().UTC(),
		Description: description_str_transaction,
		Users_Id:    &user.Id,
		Manager_Id:  *manager_id,
	}

	err = repositoryControllers.CreateManagerTransaction(manager_transaction, dbPool)
	if err != nil {
		http.Error(w, "Error Creating Manager Transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Rent Paied Sucessfully !"})
}

func GetRentById(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	rent_id_str := r.URL.Query().Get("id")

	rent_id, err := strconv.Atoi(rent_id_str)
	if err != nil {
		http.Error(w, "Invalid Rent ID", http.StatusBadRequest)
		return
	}

	rent, err := repositoryControllers.GetRentById(rent_id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching Rent", http.StatusInternalServerError)
		return
	}

	rentJson := controllers_models.Rent{
		Id:           rent.Id,
		Month:        rent.Month,
		Year:         rent.Year,
		Base_Amount:  rent.Base_Amount,
		Reduction:    rent.Reduction,
		Final_Amount: rent.Final_Amount,
		Apartment_Id: rent.Apartment_Id,
		Status:       rent.Status,
		Due_Day:      rent.Due_day,
		Pay_Day:      rent.Pay_Day,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rentJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
