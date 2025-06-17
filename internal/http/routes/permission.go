package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	handlers "github.com/sagorsarker04/Developer-Assignment/internal/http/handlers/permission"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

func RegisterPermissionRoutes(router *mux.Router) {
	// Permission Routes
	permissions := api.PathPrefix("/permissions").Subrouter()
	permissions.Use(middleware.AuthMiddleware)
	permissions.HandleFunc("/create", handlers.CreatePermission).Methods(http.MethodPost)
	
	permissions.Handle("", middleware.RequireAnyPermission([]string{"permission:read"}, http.HandlerFunc(handlers.ListAllPermissions))).Methods(http.MethodGet)                  // Admin+
	permissions.HandleFunc("/{permission_id}", handlers.GetPermissionDetails).Methods(http.MethodGet) // Admin+

	// Current User Routes
	me := api.PathPrefix("/me").Subrouter()
	me.Use(middleware.AuthMiddleware)
	me.HandleFunc("", handlers.GetCurrentUserProfile).Methods(http.MethodGet) // Authenticated
	me.HandleFunc("/permissions", handlers.GetCurrentUserPermissions).Methods(http.MethodGet) // Authenticated
}
