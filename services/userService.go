package services

import (
	"auth-api/helpers"
	"auth-api/methods"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

type UserService struct{}

// ApiRegisterUser /**
func (s *UserService) ApiRegisterUser(w http.ResponseWriter, r *http.Request) {

	// Get logged-in user data if we have it
	ctx := r.Context()
	ctxUser, getUserOk := ctx.Value("userData").(methods.User)
	if !getUserOk {
		http.Error(w, `{"error": "Could not decode user from context"}`, http.StatusInternalServerError)
		return
	}

	var newUserData = methods.CoreUserData{}

	// Decode the JSON payload
	jsonErr := json.NewDecoder(r.Body).Decode(&newUserData)
	if jsonErr != nil {
		http.Error(w, `{"error": "Invalid JSON payload"}`, http.StatusBadRequest)
		return
	}

	// Validate the data
	if newUserData.FirstName == "" {
		http.Error(w, "first_name is required", http.StatusBadRequest)
		return
	}
	if newUserData.LastName == "" {
		http.Error(w, "last_name is required", http.StatusBadRequest)
		return
	}
	if newUserData.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}
	if newUserData.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	// Sanitize Role
	validRoles := []int{1, 2} // 1 = default & 2 = admin
	if !helpers.SliceContainsInt(validRoles, newUserData.Role) {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	// Only logged in admin users can create other admins
	// if newUser role above default, check we have a logged-in user with role of at least 2
	if newUserData.Role > 1 && ctxUser.Role < 2 {
		http.Error(w, "Admin privileges required to create user with this access level", http.StatusBadRequest)
		return
	}

	// Register user in DB
	userErr := newUserData.WriteUser()
	if userErr != nil {
		http.Error(w, userErr.Error(), http.StatusBadRequest)
		return
	}

	// User created
	_, respErr := w.Write([]byte(`{"message": "User created"}`))
	if respErr != nil {
		fmt.Println(respErr.Error())
		return
	}
}

// ApiLoginUser /**
func (s *UserService) ApiLoginUser(w http.ResponseWriter, r *http.Request) {

	var loginData = methods.User{}

	// Decode the JSON payload
	jsonErr := json.NewDecoder(r.Body).Decode(&loginData)
	if jsonErr != nil {
		http.Error(w, `{"error": "Invalid JSON payload"}`, http.StatusBadRequest)
		return
	}

	// Validate the data
	if loginData.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}
	if loginData.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	// Get user by Email
	userData, dbErr := methods.GetUserBy("email", loginData.Email)
	if dbErr != nil {
		if dbErr == pgx.ErrNoRows {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		http.Error(w, dbErr.Error(), http.StatusBadRequest)
		return
	}

	// Check password
	if !methods.CheckPasswordValid(userData.Password, loginData.Password) {
		http.Error(w, "Password invalid", http.StatusBadRequest)
		return
	}

	// Generate token & update in DB
	token, tokenErr := userData.GenerateToken()
	if tokenErr != nil {
		http.Error(w, `{"error": "`+tokenErr.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Response
	jsonData, _ := json.Marshal(userData.Strip()) // encode struct as json
	_, respErr := w.Write([]byte(`{"message": "User token generated", "token" : "` + token + `", "user": ` + string(jsonData) + `}`))
	if respErr != nil {
		fmt.Println(respErr.Error())
		return
	}

}

func (s *UserService) ApiGetUser(w http.ResponseWriter, r *http.Request) {

	// Get the logged-in user
	ctx := r.Context()
	ctxUser, getUserOk := ctx.Value("userData").(methods.User)
	if !getUserOk {
		http.Error(w, `{"error": "Could not decode user from context"}`, http.StatusInternalServerError)
		return
	}

	// Get user ID from URL param if present
	userID := chi.URLParam(r, "userID")
	userIdInt, idErr := strconv.Atoi(userID)

	// The data to send back
	var userDataJson []byte

	// If fetching by ID: /api/v1/user/22
	if len(userID) > 0 {

		if idErr != nil {
			http.Error(w, "User ID must be a number", http.StatusBadRequest)
			return
		}

		// Prevent non admin fetching other users
		/** OPTIONAL DEPENDING ON REQUIREMENTS **/
		if userIdInt != ctxUser.ID && ctxUser.Role < 2 {
			http.Error(w, "Cannot fetch other profiles", http.StatusBadRequest)
			return
		}

		// Get user from DB
		dbUser, dbErr := methods.GetUserBy("id", userID)
		if dbErr != nil {
			if dbErr == pgx.ErrNoRows {
				http.Error(w, "User not found", http.StatusBadRequest)
				return
			}
			http.Error(w, dbErr.Error(), http.StatusBadRequest)
			return
		}
		userDataJson, _ = json.Marshal(dbUser.Strip()) // encode struct as json
	} else {
		// If getting logged in user: /api/v1/user
		userDataJson, _ = json.Marshal(ctxUser.Strip()) // encode struct as json
	}

	_, respErr := w.Write([]byte(`{"message": "User profile fetched", "user": ` + string(userDataJson) + `}`))
	if respErr != nil {
		fmt.Println(respErr.Error())
		return
	}
	return
}

func (s *UserService) ApiUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get logged in user
	ctx := r.Context()
	ctxUser, getUserOk := ctx.Value("userData").(methods.User)
	if !getUserOk {
		http.Error(w, `{"error": "Could not decode user from context"}`, http.StatusInternalServerError)
		return
	}

	// Update user data
	var updateUserData = methods.User{}

	// Decode the JSON payload
	jsonErr := json.NewDecoder(r.Body).Decode(&updateUserData)
	if jsonErr != nil {
		http.Error(w, `{"error": "Invalid JSON payload"}`, http.StatusBadRequest)
		return
	}

	// Validate the data
	if updateUserData.FirstName == "" {
		http.Error(w, "first_name is required", http.StatusBadRequest)
		return
	}
	if updateUserData.LastName == "" {
		http.Error(w, "last_name is required", http.StatusBadRequest)
		return
	}
	if updateUserData.ID == 0 {
		http.Error(w, "user `id` is required", http.StatusBadRequest)
		return
	}

	// Default users cannot change role to admin
	if updateUserData.Role > 1 && ctxUser.Role < 2 {
		http.Error(w, "Cannot change role", http.StatusBadRequest)
		return
	}

	// Only admin users can modify other profiles
	if ctxUser.ID != updateUserData.ID && ctxUser.Role < 2 {
		http.Error(w, "Cannot modify other profiles", http.StatusBadRequest)
		return
	}

	// Do user update
	updtErr := updateUserData.UpdateUser()
	if updtErr != nil {
		fmt.Println(updtErr.Error())
		http.Error(w, updtErr.Error(), http.StatusInternalServerError)
		return
	}

	_, respErr := w.Write([]byte(`{"message": "User updated"}`))
	if respErr != nil {
		fmt.Println(respErr.Error())
		return
	}
	return
}

func (s *UserService) ApiDeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get logged in user
	ctx := r.Context()
	ctxUser, getUserOk := ctx.Value("userData").(methods.User)
	if !getUserOk {
		http.Error(w, `{"error": "Could not decode user from context"}`, http.StatusInternalServerError)
		return
	}

	// Delete user data
	var userToDelete = methods.User{}

	// Decode the JSON payload
	jsonErr := json.NewDecoder(r.Body).Decode(&userToDelete)
	if jsonErr != nil {
		http.Error(w, `{"error": "Invalid JSON payload"}`, http.StatusBadRequest)
		return
	}

	// Validate the data
	if userToDelete.ID == 0 {
		http.Error(w, `{"user 'id' is required"}`, http.StatusBadRequest)
		return
	}

	// Only admin users can delete other profiles
	if ctxUser.ID != userToDelete.ID && ctxUser.Role < 2 {
		http.Error(w, "Cannot delete other profiles", http.StatusBadRequest)
		return
	}

	// Do user update
	updtErr := userToDelete.DeleteUser()
	if updtErr != nil {
		fmt.Println(updtErr.Error())
		http.Error(w, updtErr.Error(), http.StatusInternalServerError)
		return
	}

	_, respErr := w.Write([]byte(`{"message": "User deleted"}`))
	if respErr != nil {
		fmt.Println(respErr.Error())
		return
	}
	return
}
