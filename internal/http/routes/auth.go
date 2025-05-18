package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/handlers"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

// RegisterAuthRoutes registers the authentication-related routes.
func RegisterRoutes(router *mux.Router) {
	api := router.PathPrefix("/api/v1").Subrouter()

	// Authentication Routes
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", handlers.LoginUser).Methods(http.MethodPost)
	auth.HandleFunc("/logout", handlers.LogoutUser).Methods(http.MethodPost)
	auth.HandleFunc("/register", handlers.RegisterUser).Methods(http.MethodPost)
	auth.HandleFunc("/verify/{token}", handlers.VerifyEmail).Methods(http.MethodGet)
	auth.HandleFunc("/resend-verification", handlers.ResendVerificationEmail).Methods(http.MethodPost)
	auth.HandleFunc("/password-reset-request", handlers.PasswordResetRequest).Methods(http.MethodPost)
	auth.HandleFunc("/password-reset-confirm", handlers.PasswordResetConfirm).Methods(http.MethodPost)

	// User Routes
	users := api.PathPrefix("/users").Subrouter()
	users.Use(middleware.AuthMiddleware) // After this everyone will go through the auth middleware

	users.Handle("", middleware.RequireAnyPermission([]string{"user:read:all"}, http.HandlerFunc(handlers.ListAllUsers))).Methods(http.MethodGet)

	users.Handle("/{user_id}", middleware.RequireAnyPermission([]string{"user:read:self"}, http.HandlerFunc(handlers.GetUserDetails))).Methods(http.MethodGet)

	//ekhon kaj kore
	users.Handle("/{user_id}", middleware.RequireAnyPermission([]string{"user:update:self","user:update:all"}, http.HandlerFunc(handlers.UpdateUser))).Methods(http.MethodPut)

	users.Handle("/{user_id}",middleware.RequireAnyPermission([]string {"user:delete:self"},http.HandlerFunc(handlers.DeleteRequest))).Methods(http.MethodPost)

	users.Handle("/{user_id}", middleware.RequireAnyPermission([]string{"user:delete:all"}, http.HandlerFunc(handlers.DeleteUser))).Methods(http.MethodDelete)

	users.HandleFunc("/{user_id}/role", handlers.ChangeUserRole).Methods(http.MethodPost) // Admin+

	users.Handle("/{user_id}/promote/admin", middleware.RequireAnyPermission([]string{"user:promote:admin"}, http.HandlerFunc(handlers.PromoteToAdmin))).Methods(http.MethodPost)

	users.Handle("/{user_id}/promote/moderator", middleware.RequireAnyPermission([]string{"user:promote:moderator"}, http.HandlerFunc(handlers.PromoteToModerator))).Methods(http.MethodPost)
	
	users.Handle("/{user_id}/demote", middleware.RequireAnyPermission([]string{"user:demote"}, http.HandlerFunc(handlers.DemoteUserRole))).Methods(http.MethodPost)

	// users.HandleFunc("/{user_id}/demote", handlers.DemoteUserRole).Methods(http.MethodPost) // Admin+

	// Role Routes
	roles := api.PathPrefix("/roles").Subrouter()
	roles.Use(middleware.AuthMiddleware) // After this everyone will go through the auth middleware
	roles.Handle("", middleware.RequireAnyPermission(
		[]string{"role:read", "admin:read", "system_admin:read"},
		http.HandlerFunc(handlers.GetAllRole),
	)).Methods(http.MethodGet)

	roles.HandleFunc("/{role_id}", handlers.GetRoleDetails).Methods(http.MethodGet)
	roles.HandleFunc("/create", handlers.CreateRole).Methods(http.MethodPost)      // Admin+
	roles.HandleFunc("/{role_id}", handlers.UpdateRole).Methods(http.MethodPut)    //Only admin can update
	roles.HandleFunc("/{role_id}", handlers.DeleteRole).Methods(http.MethodDelete) // Admin+

	// Permission Routes
	permissions := api.PathPrefix("/permissions").Subrouter()
	permissions.Use(middleware.AuthMiddleware)
	permissions.HandleFunc("/create", handlers.CreatePermission).Methods(http.MethodPost)
	permissions.HandleFunc("", handlers.ListAllPermissions).Methods(http.MethodGet)                   // Admin+
	permissions.HandleFunc("/{permission_id}", handlers.GetPermissionDetails).Methods(http.MethodGet) // Admin+

	// Current User Routes
	me := api.PathPrefix("/me").Subrouter()
	me.Use(middleware.AuthMiddleware)
	me.HandleFunc("", handlers.GetCurrentUserProfile).Methods(http.MethodGet)                 // Authenticated
	me.HandleFunc("/permissions", handlers.GetCurrentUserPermissions).Methods(http.MethodGet) // Authenticated
}
