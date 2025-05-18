package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	handlers "github.com/sagorsarker04/Developer-Assignment/internal/http/handlers/auth"
)

// RegisterAuthRoutes registers the authentication-related routes.
func RegisterRoutes(router *mux.Router) {

	// Authentication Routes
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", handlers.LoginUser).Methods(http.MethodPost)
	auth.HandleFunc("/logout", handlers.LogoutUser).Methods(http.MethodPost)
	auth.HandleFunc("/register", handlers.RegisterUser).Methods(http.MethodPost)
	auth.HandleFunc("/verify/{token}", handlers.VerifyEmail).Methods(http.MethodGet)
	auth.HandleFunc("/resend-verification", handlers.ResendVerificationEmail).Methods(http.MethodPost)
	auth.HandleFunc("/password-reset-request", handlers.PasswordResetRequest).Methods(http.MethodPost)
	auth.HandleFunc("/password-reset-confirm", handlers.PasswordResetConfirm).Methods(http.MethodPost)

}
