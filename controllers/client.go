package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/portilho13/neighborconnect-backend/utils"
)

func RegisterClient(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var client controllers_models.UserCreation
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}
	// Verificar duplicação
	existingClient, err := repositoryControllers.GetUserByEmail(client.Email, dbPool)
	if err == nil && existingClient.Email != "" {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}
	encodedPassword, err := utils.GenerateFromPassword(client.Password, utils.DefaultArgon2Params)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	var apartmentID *int
	if client.ApartmentID == 0 { // If apartment id is 0 it means user as not an a
		// partment id set yet so apartment id needs to be converted to pointer to not blow the db
		apartmentID = nil
	} else {
		apartmentID = &client.ApartmentID
	}

	dbClient := models.User{
		Name:            client.Name,
		Email:           client.Email,
		Password:        encodedPassword,
		Phone:           client.Phone,
		Apartment_id:    apartmentID,
		Profile_Picture: &client.Profile_Picture,
	}

	err = repositoryControllers.CreateUser(dbClient, dbPool)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	var clientEmail string = dbClient.Email

	dbClient, err = repositoryControllers.GetUserByEmail(clientEmail, dbPool)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	userAccount := models.Account{
		Account_number: utils.GenerateRandomHash(),
		Balance:        0,
		Currency:       "EUR",
		Users_id:       &dbClient.Id,
	}

	err = repositoryControllers.CreateAccount(userAccount, dbPool)
	if err != nil {
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return
	}

	err = repositoryControllers.UpdateApartmentStatus(*apartmentID, dbPool)
	if err != nil {
		http.Error(w, "Error updating apartment status", http.StatusInternalServerError)
		return
	}

	manager_id, err := repositoryControllers.GetManagerIdByApartmentId(*apartmentID, dbPool)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error fetching manager", http.StatusInternalServerError)
		return
	}

	manager_activity_str := fmt.Sprintf("%s registered to apartment %d", client.Name, *apartmentID)
	manager_activity := models.Manager_Activity{
		Type:        "New Resident",
		Description: manager_activity_str,
		Created_At:  time.Now().UTC(),
		Manager_Id:  *manager_id,
	}

	err = repositoryControllers.CreateManagerActivity(manager_activity, dbPool)
	if err != nil {
		http.Error(w, "Error creating manager activity", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Client Registed !"})
}

func LoginClient(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {

	// var creds controllers_models.UserJson
	var creds controllers_models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch user from the database
	user, err := repositoryControllers.GetUserByEmail(creds.Email, dbPool)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password
	match, err := utils.ComparePasswordAndHash(creds.Password, user.Password)
	if err != nil || !match {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	// Create a session
	session, _ := utils.Store.Get(r, "session")
	session.Values["user_id"] = user.Id
	session.Values["email"] = user.Email
	session.Values["role"] = "user"
	session.Save(r, w)

	avatar := ""
	if user.Profile_Picture != nil {
		avatar = *user.Profile_Picture
	}

	userJson := controllers_models.UserLogin{
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		Phone:       user.Phone,
		ApartmentID: *user.Apartment_id,
		Avatar:      avatar,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userJson)
}

func LogoutHandlerUser(w http.ResponseWriter, r *http.Request) {

	session, err := utils.Store.Get(r, "session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	// Limpar os dados da sessão
	delete(session.Values, "user_id")
	delete(session.Values, "email")
	delete(session.Values, "role")

	// Invalidar sessão
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}
	role := session.Values["role"]
	fmt.Println("Logging out role:", role)
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1, // Força a expiração
		HttpOnly: true,
	}
	session.Save(r, w)

	log.Println(session.Values)
	// Enviar resposta ao cliente
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}

func GetNeighborInfo(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	users, err := repositoryControllers.GetAllUsers(dbPool)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	var usersJson []controllers_models.NeighborInfo

	for _, user := range users {
		userJson := controllers_models.NeighborInfo{
			Id:              user.Id,
			Name:            user.Name,
			Email:           user.Email,
			Phone:           user.Phone,
			ApartmentID:     *user.Apartment_id,
			Profile_Picture: *user.Profile_Picture,
		}

		usersJson = append(usersJson, userJson)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(usersJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func UploadProfilePicture(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	user_id_str := r.FormValue("user_id")
	user_id, err := strconv.Atoi(user_id_str)
	if err != nil {
		http.Error(w, "Failed to convert user id", http.StatusInternalServerError)
		return
	}

	var api_path string

	files := r.MultipartForm.File["profilePicture"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to create user file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		filename := fileHeader.Filename
		ext := filepath.Ext(filename)
		newFilename := uuid.New().String() + ext
		savePath := fmt.Sprintf("./uploads/users/%s", newFilename)

		dst, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "Failed to create user file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to create listing", http.StatusInternalServerError)
			return
		}

		api_url := utils.GetApiUrl()

		api_path = fmt.Sprintf("https://%s/api/v1/uploads/users/%s", api_url, newFilename)

	}

	err = repositoryControllers.UpdateUserProfilePicture(api_path, user_id, dbPool)
	if err != nil {
		http.Error(w, "Failed to update user profile picture", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Client Profile Picture Updated !"})
}
