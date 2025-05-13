package main

import (
	"log"
	"os"

	"github.com/sagorsarker04/Developer-Assignment/internal/database"
)

func main() {
	// Set the default migrations folder
	migrationsPath := "./migrations"
	if len(os.Args) > 1 {
		migrationsPath = os.Args[1]
	}

	// Run migrations
	if err := database.Migrate(migrationsPath); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully.")
}
