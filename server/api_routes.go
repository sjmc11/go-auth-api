package server

import (
	"auth-api/methods"
	"auth-api/services"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"net/http"
)

func initApiRoutes(r *chi.Mux) error {

	appService := services.Init()

	r.Route("/api/v1", func(r chi.Router) {

		// Set JSON response type for all API routes
		r.Use(middleware.SetHeader("Content-Type", "application/json"))

		/**
		PUBLIC ROUTES
		*/

		// Heartbeat
		r.Get("/", appService.ApiHeartBeat)

		// Register
		r.With(tokenCheckMiddleware(0)).Post("/register", appService.UserService.ApiRegisterUser)

		// Login
		r.Post("/login", appService.UserService.ApiLoginUser)

		// Request reset password email
		r.Post("/reset-password/request", appService.AuthService.ApiPasswordResetEmail)

		// Check reset password auth code
		r.Post("/reset-password/check", appService.AuthService.ApiCheckPwResetCode)

		// Do password reset
		r.Post("/reset-password/update", appService.AuthService.ApiDoPasswordReset)

		/**
		PRIVATE ROUTES
		*/

		// Get logged in user
		r.With(tokenCheckMiddleware(1)).Get("/user", appService.UserService.ApiGetUser)

		// Get user by ID
		r.With(tokenCheckMiddleware(1)).Get("/user/{userID}", appService.UserService.ApiGetUser)

		// User update
		r.With(tokenCheckMiddleware(1)).Post("/user/update", appService.UserService.ApiUpdateUser)

		// User delete
		r.With(tokenCheckMiddleware(1)).Post("/user/delete", appService.UserService.ApiDeleteUser)

		// Admin only - LEVEL 2+
		r.With(tokenCheckMiddleware(2)).Get("/admin", appService.ApiHeartBeat)

	})

	return nil
}

/**
tokenCheckMiddleware
Access level 0 = All users
Access level 1 = Logged-in users
Access level 2 = Admin only
*/
var tokenCheckMiddleware = func(accessLevel int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Get token
			headerToken := r.Header.Get("token")
			if headerToken == "" && accessLevel > 0 {
				http.Error(w, `{"error": "No token detected"}`, http.StatusUnauthorized)
				return
			}

			var userData methods.User
			var dbErr error

			// Get user by Token
			if len(headerToken) > 0 {
				userData, dbErr = methods.GetUserBy("token", headerToken)
				if dbErr != nil {
					if dbErr == pgx.ErrNoRows {
						http.Error(w, `{"error": "Token invalid"}`, http.StatusBadRequest)
						return
					}
					http.Error(w, dbErr.Error(), http.StatusBadRequest)
					return
				}
			}

			// Access level/Role check
			if accessLevel > 0 && userData.Role < accessLevel {
				http.Error(w, `{"error": "Permission denied"}`, http.StatusForbidden)
				return
			}

			fmt.Printf("%+v\n", userData)

			ctx := context.WithValue(r.Context(), "userData", userData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
