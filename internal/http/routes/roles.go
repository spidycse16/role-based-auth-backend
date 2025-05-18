package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	handlers "github.com/sagorsarker04/Developer-Assignment/internal/http/handlers/role"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

func RoleRoutes(router *mux.Router) {
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

	roles.HandleFunc("/{user_id}/role", handlers.ChangeUserRole).Methods(http.MethodPost) // Admin+

	roles.Handle("/{user_id}/promote/admin", middleware.RequireAnyPermission([]string{"user:promote:admin"}, http.HandlerFunc(handlers.PromoteToAdmin))).Methods(http.MethodPost)

	roles.Handle("/{user_id}/promote/moderator", middleware.RequireAnyPermission([]string{"user:promote:moderator"}, http.HandlerFunc(handlers.PromoteToModerator))).Methods(http.MethodPost)

	roles.Handle("/{user_id}/demote", middleware.RequireAnyPermission([]string{"user:demote"}, http.HandlerFunc(handlers.DemoteUserRole))).Methods(http.MethodPost)

}
