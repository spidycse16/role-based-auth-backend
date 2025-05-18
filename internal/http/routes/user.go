package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	handlers "github.com/sagorsarker04/Developer-Assignment/internal/http/handlers/user"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

func RegisterUserRoutes(r *mux.Router) {
	// User Routes
	users := api.PathPrefix("/users").Subrouter()
	users.Use(middleware.AuthMiddleware) // After this everyone will go through the auth middleware

	users.Handle("", middleware.RequireAnyPermission([]string{"user:read:all"}, http.HandlerFunc(handlers.ListAllUsers))).Methods(http.MethodGet)

	users.Handle("/{user_id}", middleware.RequireAnyPermission([]string{"user:read:self"}, http.HandlerFunc(handlers.GetUserDetails))).Methods(http.MethodGet)

	//ekhon kaj kore
	users.Handle("/{user_id}", middleware.RequireAnyPermission([]string{"user:update:self", "user:update:all"}, http.HandlerFunc(handlers.UpdateUser))).Methods(http.MethodPut)

	users.Handle("/{user_id}", middleware.RequireAnyPermission([]string{"user:delete:self"}, http.HandlerFunc(handlers.DeleteRequest))).Methods(http.MethodPost)

	users.Handle("/{user_id}", middleware.RequireAnyPermission([]string{"user:delete:all"}, http.HandlerFunc(handlers.DeleteUser))).Methods(http.MethodDelete)

	
	// users.HandleFunc("/{user_id}/demote", handlers.DemoteUserRole).Methods(http.MethodPost) // Admin+
}
