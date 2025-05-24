package main

import (
	"log"

	"github.com/sagorsarker04/Developer-Assignment/internal/config"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
)


func main() {
	// Load configuration
	cfg := config.GetConfig()
	log.Print("Chill",cfg)
	// connect to database
	database.Connect()
	defer database.Close()
	// add system admin if does not exist
	database.InitAdminUser(cfg.Admin)
}
