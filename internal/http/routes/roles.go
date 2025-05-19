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
		[]string{"role:read"},
		http.HandlerFunc(handlers.GetAllRole),
	)).Methods(http.MethodGet)

	roles.Handle("/{role_id}", middleware.RequireAnyPermission([]string{"role:read"}, http.HandlerFunc(handlers.GetRoleDetails))).Methods(http.MethodGet)

	roles.Handle("/create", middleware.RequireAnyPermission([]string{"role:create"}, http.HandlerFunc(handlers.CreateRole))).Methods(http.MethodPost)

	roles.Handle("/{role_id}", middleware.RequireAnyPermission([]string{"role:update"}, http.HandlerFunc(handlers.UpdateRole))).Methods(http.MethodPut)

	roles.Handle("/{role_id}", middleware.RequireAnyPermission([]string{"role:delete"}, http.HandlerFunc(handlers.DeleteRole))).Methods(http.MethodDelete)
	
	roles.Handle("/{user_id}/role", middleware.RequireAnyPermission([]string{"role:update","user:update:all"}, http.HandlerFunc(handlers.ChangeUserRole))).Methods(http.MethodPost)

	roles.Handle("/{user_id}/promote/admin", middleware.RequireAnyPermission([]string{"user:promote:admin"}, http.HandlerFunc(handlers.PromoteToAdmin))).Methods(http.MethodPost)

	roles.Handle("/{user_id}/promote/moderator", middleware.RequireAnyPermission([]string{"user:promote:moderator"}, http.HandlerFunc(handlers.PromoteToModerator))).Methods(http.MethodPost)

	roles.Handle("/{user_id}/demote", middleware.RequireAnyPermission([]string{"user:demote"}, http.HandlerFunc(handlers.DemoteUserRole))).Methods(http.MethodPost)

}
