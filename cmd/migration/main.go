package main

import (
	"fmt"


	"github.com/sagorsarker04/Developer-Assignment/internal/config"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
)


func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		 fmt.Println("Failed to load config")
		 return
	}
	// connect to database
	database.Connect()

	// add system admin if does not exist
	database.InitAdminUser(cfg.Admin)
}
