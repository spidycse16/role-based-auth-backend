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
	// auth.HandleFunc("/resend-verification", handlers.ResendVerificationEmail).Methods(http.MethodPost)
	// auth.HandleFunc("/password-reset", handlers.ResetPassword).Methods(http.MethodPost)

	// // User Routes
	// users := api.PathPrefix("/users").Subrouter()
	// users.HandleFunc("", handlers.ListAllUsers).Methods(http.MethodGet) // Admin+
	// users.HandleFunc("/{user_id}", handlers.GetUserDetails).Methods(http.MethodGet) // User(self only)/Moderator+
	// users.HandleFunc("/{user_id}", handlers.UpdateUser).Methods(http.MethodPut) // User(self only)/Admin+
	// users.HandleFunc("/{user_id}/request-deletion", handlers.RequestAccountDeletion).Methods(http.MethodPost) // User(self only)
	// users.HandleFunc("/{user_id}", handlers.DeleteUser).Methods(http.MethodDelete) // Moderator+
	// users.HandleFunc("/{user_id}/role", handlers.ChangeUserRole).Methods(http.MethodPost) // Admin+
	// users.HandleFunc("/{user_id}/promote/admin", handlers.PromoteToAdmin).Methods(http.MethodPost) // System Admin
	// users.HandleFunc("/{user_id}/promote/moderator", handlers.PromoteToModerator).Methods(http.MethodPost) // Admin+
	// users.HandleFunc("/{user_id}/demote", handlers.DemoteUserRole).Methods(http.MethodPost) // Admin+

	// Role Routes
	roles := api.PathPrefix("/roles").Subrouter()
	roles.Use(middleware.AuthMiddleware) // After this everyone will go through the auth middleware
	roles.HandleFunc("", handlers.GetAllRole).Methods(http.MethodGet)
	roles.HandleFunc("/{role_id}", handlers.GetRoleDetails).Methods(http.MethodGet)
	roles.HandleFunc("/create", handlers.CreateRole).Methods(http.MethodPost) // Admin+
	roles.HandleFunc("/{role_id}", handlers.UpdateRole).Methods(http.MethodPut) //Only admin can update
	roles.HandleFunc("/{role_id}", handlers.DeleteRole).Methods(http.MethodDelete) // Admin+

	// Permission Routes
	permissions := api.PathPrefix("/permissions").Subrouter()
    permissions.Use(middleware.AuthMiddleware)
    permissions.HandleFunc("/create", handlers.CreatePermission).Methods(http.MethodPost)
	// permissions.HandleFunc("", handlers.ListAllPermissions).Methods(http.MethodGet) // Admin+
	// permissions.HandleFunc("/{permission_id}", handlers.GetPermissionDetails).Methods(http.MethodGet) // Admin+

	// Current User Routes
	me := api.PathPrefix("/me").Subrouter()
    me.Use(middleware.AuthMiddleware)
	me.HandleFunc("", handlers.GetCurrentUserProfile).Methods(http.MethodGet) // Authenticated
	me.HandleFunc("/permissions", handlers.GetCurrentUserPermissions).Methods(http.MethodGet) // Authenticated
}
