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
	"github.com/sagorsarker04/Developer-Assignment/internal/config"

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

	sqlFilePath := "/app/migrations/000001_init_schema/up.sql"
	//sqlFilePath := "/home/shahid/Desktop/Developer-Assignment/migrations/000001_init_schema/up.sql"

	sqlBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		log.Println("failed to read sqlFilePath")
	}

	_, err = DB.Exec(string(sqlBytes))

	if err != nil {
		log.Println("DB execute error")
		return
	}
	log.Println("Hurray, database created..")
	var exists bool
	err = DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)", admin.Email).Scan(&exists)
	if err != nil {
		log.Fatal("DB check failed:", err)
	}

	// salt_pass := os.Getenv("PASSWORD_SALT")
	// if !exists {
	// 	hashed, err := bcrypt.GenerateFromPassword([]byte(admin.Password+salt_pass), bcrypt.DefaultCost)
	// 	if err != nil {
	// 		log.Fatal("Hash error:", err)
	// 	}

	// Insert user and get the user_id
	// var user_id uuid.UUID
	// err = DB.QueryRow(`
	//     INSERT INTO users (username, email, password_hash, user_type)
	//     VALUES ($1, $2, $3, $4)
	//     RETURNING id`,
	// 	admin.Username, admin.Email, string(hashed), admin.UserType).Scan(&user_id)

	// if err != nil {
	// 	log.Fatal("Admin insert failed:", err)
	// }

	// Get role_id for 'user' role
	// var role_id uuid.UUID
	// err = DB.QueryRow(`
	//     SELECT id FROM roles
	//     WHERE name = $1`,
	// 	admin.UserType).Scan(&role_id)

	// if err != nil {
	// 	log.Fatal("Failed to get role id:", err)
	// }

	// Assign role to user
	// 	_, err = DB.Exec(`
	//         INSERT INTO user_roles (user_id, role_id, assigned_by)
	//         VALUES ($1, $2, $1)`,
	// 		user_id, role_id)

	// 	if err != nil {
	// 		log.Fatal("Failed to assign role to admin:", err)
	// 	}

	// 	fmt.Println("System admin registered.")
	// } else {
	// 	fmt.Println("System admin already exists.")
	// }
}
