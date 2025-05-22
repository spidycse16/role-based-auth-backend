package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/sagorsarker04/Developer-Assignment/internal/config"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect establishes a connection to the PostgreSQL database using configuration.
func Connect() (*sql.DB, error) {
	// Load configuration singletone varible
	cfg := config.GetConfig()
	fmt.Printf("Address of cfg 2: %p\n", cfg)
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
	db, err := sql.Open(cfg.Database.Host, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	DB = db
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

func InitAdminUser(admin config.AdminConfig) {
	// Path to migration file
	sqlFilePath := "/app/migrations/000001_init_schema/up.sql" // Adjust for your Docker setup

	// Read the SQL file
	sqlBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	// Execute the migration SQL
	_, err = DB.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("DB schema execution error: %v", err)
	}
	log.Println("✅ Database schema initialized.")

	// Check if admin already exists
	var exists bool
	err = DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", admin.Email).Scan(&exists)
	if err != nil {
		log.Fatalf("Error checking for existing admin user: %v", err)
	}

	if exists {
		log.Println("✅ System admin already exists.")
		return
	}

	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Password hashing failed: %v", err)
	}

	// Insert admin user and get user ID
	var userID uuid.UUID
	err = DB.QueryRow(`
		INSERT INTO users (username, email, password_hash, user_type, email_verified)
		VALUES ($1, $2, $3, $4, TRUE)
		RETURNING id`,
		admin.Username, admin.Email, string(hashedPassword), "system_admin").Scan(&userID)

	if err != nil {
		log.Fatalf("Failed to insert system admin user: %v", err)
	}

	// Get role ID for 'system_admin'
	var roleID uuid.UUID
	err = DB.QueryRow(`SELECT id FROM roles WHERE name = 'system_admin'`).Scan(&roleID)
	if err != nil {
		log.Fatalf("Failed to get role ID for system_admin: %v", err)
	}

	// Assign role to the user (assigning to self)
	_, err = DB.Exec(`
		INSERT INTO user_roles (user_id, role_id, assigned_by)
		VALUES ($1, $2, $1)`,
		userID, roleID)

	if err != nil {
		log.Fatalf("Failed to assign role to system admin: %v", err)
	}

	log.Println("🎉 System admin user created and role assigned.")
}
