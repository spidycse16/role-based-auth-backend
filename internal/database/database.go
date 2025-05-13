package database

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sagorsarker04/Developer-Assignment/internal/config"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Connect establishes a connection to the PostgreSQL database using configuration.
func Connect() (*sql.DB, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Build connection string
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	// Open database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully.")
	return db, nil
}

// Close terminates the database connection.
func Close(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}
}

// Migrate runs all the database migrations.
func Migrate(migrationsPath string) error {
	db, err := Connect()
	if err != nil {
		return err
	}
	defer Close(db)

	// Get the database driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL driver: %w", err)
	}

	// Build the migrations source path
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to resolve migrations path: %w", err)
	}

	// Create the migration instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Run the migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Database migrations applied successfully.")
	return nil
}
