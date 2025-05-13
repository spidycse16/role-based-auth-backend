package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/routes"
)

func main() {
	router := mux.NewRouter()

	// Register auth routes
	routes.RegisterAuthRoutes(router)

	// Start the server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
