package routes

import "github.com/gorilla/mux"

var api *mux.Router

func SetupRoutes(router *mux.Router) {
	api = router.PathPrefix("/api/v1").Subrouter()
	RegisterRoutes(router)
	RoleRoutes(router)
	RegisterPermissionRoutes(router)
	RegisterUserRoutes(router)
}
