package services

import (
	"fmt"
	"net/http"
)

type AppServices struct {
	UserService *UserService
	AuthService *AuthService
}

func Init() *AppServices {
	return &AppServices{
		UserService: &UserService{},
		AuthService: &AuthService{},
	}
}

func (s *AppServices) ApiHeartBeat(w http.ResponseWriter, r *http.Request) {
	_, respErr := w.Write([]byte(`{"message": "Welcome to the GoAuth API"}`))
	if respErr != nil {
		fmt.Println(respErr.Error())
		return
	}
}
