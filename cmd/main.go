package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/config"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/routes"
)

func main() {
	router := mux.NewRouter()

	routes.SetupRoutes(router)

	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:5173"})

	allowedMethods := handlers.AllowedMethods([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
	})

	allowedHeaders := handlers.AllowedHeaders([]string{
		"Content-Type",
		"Authorization",
		"Cookie",
	})

	allowedCredentials := handlers.AllowCredentials()
	cfg := config.GetConfig()
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("Address of cfg: %p\n", cfg)
	log.Printf("Server running on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr,
		handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders, allowedCredentials)(router),
	))
}
