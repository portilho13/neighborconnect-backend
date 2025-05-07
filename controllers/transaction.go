package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/email"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	repositoryControllersUsers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"

	"github.com/portilho13/neighborconnect-backend/utils"
)

const FEES = 0.05 //5%, maybe then put in db ??

func PayTransaction(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var payJson controllers_models.PayJson

	err := json.NewDecoder(r.Body).Decode(&payJson)
	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	transaction, err := repositoryControllers.GetTransactionById(payJson.Transaction_Id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching Transaction", http.StatusInternalServerError)
		return
	}

	if transaction.Payment_Status != "pending" {
		http.Error(w, "Invalid Transaction", http.StatusInternalServerError)
		return
	}

	if *transaction.Buyer_Id != payJson.User_Id {
		http.Error(w, "Invalid User", http.StatusInternalServerError)
		return
	}

	switch payJson.Type {
	case "wallet":
		if !utils.ValidateWalletBalance(payJson.User_Id, transaction.Final_Price, dbPool) {
			http.Error(w, "Insuficient Balance", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Unsuported Method", http.StatusInternalServerError)
		return
	}

	seller_account, err := repositoryControllersUsers.GetAccountByUserId(*transaction.Seller_Id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching Seller Id Account", http.StatusInternalServerError)
		return
	}

	payout := transaction.Final_Price * (1 - FEES)

	newBalance := seller_account.Balance + payout

	err = repositoryControllersUsers.UpdateAccountBalance(seller_account.Id, newBalance, dbPool)
	if err != nil {
		http.Error(w, "Error Adding Payout", http.StatusInternalServerError)
		return
	}

	err = repositoryControllers.UpdateTransactionStatus("paid", *transaction.Id, dbPool)
	if err != nil {
		http.Error(w, "Error Changing Transaction Status", http.StatusInternalServerError)
		return
	}

	user, err := repositoryControllersUsers.GetUsersById(payJson.User_Id, dbPool)
	if err != nil {
		http.Error(w, "Error Getting User Id", http.StatusInternalServerError)
		return
	}

	email_struct := email.Email{
		To:      []string{user.Email},
		Subject: "Order Receipt",
	}

	err = email.SendEmail(email_struct, "order receipt", transaction)
	if err != nil {
		http.Error(w, "Error Sending Email", http.StatusInternalServerError)
		return
	}

	feesAmount := transaction.Final_Price * FEES

	manager_id, err := utils.GetManagerIdByUserId(*transaction.Seller_Id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching Manager Id", http.StatusInternalServerError)
		return
	}

	manager_transaction := models.Manager_Transaction{
		Type:        "fees",
		Amount:      feesAmount,
		Date:        time.Now(),
		Description: "Marketplace Fees",
		Manager_Id:  *manager_id,
	}

	err = repositoryControllersUsers.CreateManagerTransaction(manager_transaction, dbPool)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Adding Manager Transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Pay Concluded Sucessfully !"})
}

func GetAllTransactions(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	query := r.URL.Query()
	idStr := query.Get("user_id")

	user_id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	transactions, err := repositoryControllers.GetTransactionsByBuyerId(user_id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching User", http.StatusInternalServerError)
		return
	}

	var transactionsJson []controllers_models.TransactionJson

	for _, transaction := range transactions {
		transactionJson := controllers_models.TransactionJson{
			Id:               transaction.Id,
			Final_Price:      transaction.Final_Price,
			Transaction_Time: transaction.Transaction_Time,
			Transaction_Type: transaction.Transaction_Type,
			Seller_Id:        transaction.Seller_Id,
			Buyer_Id:         transaction.Buyer_Id,
			Listing_Id:       transaction.Listing_Id,
			Payment_Status:   transaction.Payment_Status,
			Payment_Due_time: transaction.Payment_Due_time,
		}

		transactionsJson = append(transactionsJson, transactionJson)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transactionsJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
