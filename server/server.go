package server

import (
	"auth-api/helpers/env"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/logrusorgru/aurora/v4"
	"net/http"
)

func ServeApp() error {
	appChi := chi.NewRouter()

	// CORE MIDDLEWARES
	appChi.Use(middleware.Logger)
	appChi.Use(middleware.RequestID)
	appChi.Use(middleware.RealIP)
	appChi.Use(middleware.Recoverer)

	// LOAD API ROUTES
	apiRoutesErr := initApiRoutes(appChi)
	if apiRoutesErr != nil {
		return apiRoutesErr
	}

	// START SERVER
	fmt.Println(aurora.Green("✓ STARTING APP - http://localhost:" + env.Get("WEBPORT")))
	err := http.ListenAndServe(":"+env.Get("WEBPORT"), appChi)
	if err != nil {
		return err
	}

	return nil
}
