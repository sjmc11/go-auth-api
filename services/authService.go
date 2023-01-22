package services

import (
	"auth-api/helpers"
	"auth-api/methods"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strings"
)

type AuthService struct{}

// ApiCheckPwResetCode /**
func (s *AuthService) ApiCheckPwResetCode(w http.ResponseWriter, r *http.Request) {

	// JSON Data structure
	var PwResetCheck = struct {
		AuthCode string `json:"auth_code"`
	}{}

	// Decode the JSON payload
	jsonErr := json.NewDecoder(r.Body).Decode(&PwResetCheck)
	if jsonErr != nil {
		http.Error(w, `{"error": "Invalid JSON payload"}`, http.StatusBadRequest)
		return
	}

	if PwResetCheck.AuthCode == "" {
		http.Error(w, "auth_code required", http.StatusBadRequest)
		return
	}

	pwReset, checkErr := methods.GetPwReset(PwResetCheck.AuthCode)
	if checkErr != nil {
		if checkErr == pgx.ErrNoRows {
			http.Error(w, "Auth code invalid", http.StatusBadRequest)
			return
		}
		http.Error(w, checkErr.Error(), http.StatusInternalServerError)
		return
	}

	authCodeValid, reason := pwReset.IsValid()
	if !authCodeValid {
		http.Error(w, reason, http.StatusBadRequest)
		return
	}

	pwResetJson, _ := json.Marshal(pwReset) // encode struct as json

	_, respErr := w.Write([]byte(`{"message": "Auth code is ok", "pw_reset": ` + string(pwResetJson) + `}`))
	if respErr != nil {
		fmt.Println(respErr.Error())
		return
	}
}

// ApiPasswordResetEmail /**
func (s *AuthService) ApiPasswordResetEmail(w http.ResponseWriter, r *http.Request) {

	// JSON Data structure
	var pwResetForm = struct {
		Email string `json:"email"`
	}{}

	// Decode the JSON payload
	jsonErr := json.NewDecoder(r.Body).Decode(&pwResetForm)
	if jsonErr != nil {
		http.Error(w, `{"error": "Invalid JSON payload"}`, http.StatusBadRequest)
		return
	}

	// Validate email address
	if !helpers.EmailValid(pwResetForm.Email) {
		http.Error(w, "Enter a valid email address", http.StatusBadRequest)
		return
	}

	// Get user data
	userData, dbErr := methods.GetUserBy("email", pwResetForm.Email)
	if dbErr != nil {
		if dbErr == pgx.ErrNoRows {
			http.Error(w, "No matching user", http.StatusBadRequest)
			return
		}
		http.Error(w, dbErr.Error(), http.StatusBadRequest)
		return
	}

	// Generate a pw reset record
	authCode, authCdErr := userData.CreatePwReset()
	if authCdErr != nil {
		http.Error(w, authCdErr.Error(), http.StatusBadRequest)
		return
	}

	// Email it to the user
	_, emailErr := userData.SendPwResetEmail(authCode)
	if emailErr != nil {
		fmt.Println(emailErr.Error())
		return
	}

	// JSON response
	_, respErr := w.Write([]byte(`{"message": "Password reset generated, check email for ` + userData.Email + `", "auth_code" : "` + authCode + `"}`))
	if respErr != nil {
		fmt.Println(respErr.Error())
		return
	}
	return
}

// ApiDoPasswordReset /**
func (s *AuthService) ApiDoPasswordReset(w http.ResponseWriter, r *http.Request) {

	// JSON Data structure
	var pwResetForm = struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		PassConfirm string `json:"pass_confirm"`
		AuthCode    string `json:"auth_code"`
	}{}

	// Decode the JSON payload
	jsonErr := json.NewDecoder(r.Body).Decode(&pwResetForm)
	if jsonErr != nil {
		http.Error(w, `{"error": "Invalid JSON payload"}`, http.StatusBadRequest)
		return
	}

	// Validate email address
	if !helpers.EmailValid(pwResetForm.Email) {
		http.Error(w, `{"error": "Enter a valid email address"}`, http.StatusBadRequest)
		return
	}

	// Get the pw_reset record
	pwReset, checkErr := methods.GetPwReset(pwResetForm.AuthCode)
	if checkErr != nil {
		if checkErr == pgx.ErrNoRows {
			http.Error(w, `{"error": "Auth code invalid"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, checkErr.Error(), http.StatusInternalServerError)
		return
	}

	// Check if already used or expired
	authCodeValid, reason := pwReset.IsValid()
	if !authCodeValid {
		http.Error(w, reason, http.StatusBadRequest)
		return
	}

	// Make sure the email matches up to the auth_code
	if helpers.SanitizeString(pwResetForm.Email) != helpers.SanitizeString(pwReset.UserEmail) {
		http.Error(w, `{"error": "This email address does not match the reset code"}`, http.StatusBadRequest)
		return
	}

	// Make sure the passwords match
	if strings.TrimSpace(pwResetForm.Password) != strings.TrimSpace(pwResetForm.PassConfirm) {
		http.Error(w, `{"error": "Passwords do not match"}`, http.StatusBadRequest)
		return
	}

	// Ephemeral user object - populate the user ID and email which is minimum requirement for most methods
	ephUserData := methods.User{
		CoreUserData: methods.CoreUserData{
			Email: pwReset.UserEmail,
		},
		ID: pwReset.UserID,
	}

	// Update password in DB
	passErr := ephUserData.UpdatePassword(strings.TrimSpace(pwResetForm.Password))
	if passErr != nil {
		http.Error(w, `{"error": "`+passErr.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Generate a fresh token for authentication
	token, tokenErr := ephUserData.GenerateToken()
	if tokenErr != nil {
		http.Error(w, `{"error": "`+tokenErr.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Mark pw_reset as 'used'
	setUsedErr := pwReset.MarkAsUsed()
	if setUsedErr != nil {
		http.Error(w, `{"error": "`+setUsedErr.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// JSON response
	_, respErr := w.Write([]byte(`{"message": "Password updated", "token": "` + token + `"}`))
	if respErr != nil {
		fmt.Println(respErr.Error())
		return
	}
	return
}
