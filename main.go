package main

import (
	"auth-api/db"
	"auth-api/helpers/env"
	"auth-api/mailer"
	"auth-api/server"
	"fmt"
	"github.com/logrusorgru/aurora/v4"
	"log"
)

func main() {

	// .ENV
	initEnv := env.Environment{EnvPath: ".env"}
	initEnv.LoadEnv()

	// DATABASE
	dbErr := db.PgConnection.InitDbConn()
	if dbErr != nil {
		log.Fatal(aurora.Red("Error Connecting to Database: " + dbErr.Error()))
		return
	}
	db.PgConnection.PingDB()

	// SETUP MAILER
	_, mailErr := mailer.InitGoMailer()
	if mailErr != nil {
		fmt.Println(mailErr.Error())
		return
	}

	// SERVE
	serverErr := server.ServeApp()
	if serverErr != nil {
		log.Fatal(aurora.Red("Error serving app: " + serverErr.Error()))
		return
	}
}
