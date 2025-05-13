package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/handlers"
    // "github.com/sagorsarker04/Developer-Assignment/internal/services/email"
)

// RegisterAuthRoutes registers the authentication-related routes.
func RegisterAuthRoutes(router *mux.Router) {
	api := router.PathPrefix("/api/v1/auth").Subrouter()
	api.HandleFunc("/register", handlers.RegisterUser).Methods(http.MethodPost)
	api.HandleFunc("/login", handlers.LoginUser).Methods(http.MethodPost)
	api.HandleFunc("/verify/{token}", handlers.VerifyEmail).Methods(http.MethodGet)
}
