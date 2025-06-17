#!/bin/bash
# Setup script for Developer Assignment Auth Service

# Create main structure
mkdir -p cmd
mkdir -p deployments
mkdir -p internal/config
mkdir -p internal/database
mkdir -p internal/models
mkdir -p internal/services
mkdir -p internal/http/handlers
mkdir -p internal/http/middleware
mkdir -p internal/http/routes
mkdir -p migrations
mkdir -p scripts
mkdir -p static
mkdir -p test

# Create .gitkeep files for empty directories
touch deployments/.gitkeep
touch internal/config/.gitkeep
touch internal/database/.gitkeep
touch internal/models/.gitkeep
touch internal/services/.gitkeep
touch internal/http/handlers/.gitkeep
touch internal/http/middleware/.gitkeep
touch internal/http/routes/.gitkeep
touch migrations/.gitkeep
touch scripts/.gitkeep
touch static/.gitkeep
touch test/.gitkeep

# Create main.go file
cat > cmd/main.go << 'EOF'
package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("AffPilot Auth Service starting...")
	log.Println("Server initialized")
}
EOF

# Create README.md file
cat > README.md << 'EOF'
# 🚀 AffPilot Auth Service

A robust authentication and authorization service for managing user identities, issuing JWT tokens, and implementing role-based access controls.

> **DEVELOPER ASSIGNMENT**: This project serves as a technical assessment for backend engineers. Candidates should implement the authentication service according to the specifications below, with particular attention to the database schema and API endpoints.

## Overview

The Authentication Service provides:
- Secure user authentication (username/password)
- JWT token issuance
- Comprehensive identity and access management (IAM)
- Email verification workflow

## Features

### Authentication
- Username/password login
- Secure password hashing with bcrypt
- Session management with JWT
- Email verification with time-limited hash links

### Authorization
- Role-based permissions system with four user types:
  - **System Admin**: Full system access (defined at initialization)
  - **Admin**: Full access except cannot modify other admins
  - **Moderator**: Can manage all data but cannot delete users
  - **User**: Can only manage their own data
- Permission-based access control

### User Management
- Create, read, update user accounts
- Role assignment (admin only)
- Email verification workflow
- User deletion request workflow

### Administration
- Initial system admin credentials defined in environment variables
- Protected admin functions
- Audit trail for role assignments
- Tiered administrative privileges

## Technology Stack

| Component    | Technology                           |
|--------------|--------------------------------------|
| Language     | Go (Golang)                          |
| Database     | PostgreSQL                           |
| Security     | JWT, bcrypt                          |
| Deployment   | Docker                               |

## Project Structure

```
.
├── cmd/
│   └── main.go
├── deployments/
├── internal/
│   ├── config/         # App configurations
│   ├── database/       # Database connections, migrations
│   ├── models/         # Data models
│   ├── services/       # Business logic
│   └── http/
│       ├── handlers/   # HTTP handlers
│       ├── middleware/ # Middleware components
│       └── routes/     # API route definitions
├── migrations/         # SQL migration files
├── scripts/            # Utility scripts
├── static/             # Static assets
├── test/               # Test cases and mocks
├── .env.example        # Environment variables template
├── .gitignore
├── README.md
├── go.mod
└── go.sum
```

## Database Schema

The following tables should be implemented in your PostgreSQL database:

### Users

| Column             | Type             | Constraints                 | Description                         |
|--------------------|------------------|----------------------------|-------------------------------------|
| id                 | uuid             | PRIMARY KEY                | Unique user identifier              |
| username           | varchar(50)      | UNIQUE, NOT NULL          | User's login name                   |
| email              | varchar(100)     | UNIQUE, NOT NULL          | User's email address                |
| password_hash      | varchar(100)     | NOT NULL                  | bcrypt hashed password              |
| first_name         | varchar(50)      | NULL                      | User's first name                   |
| last_name          | varchar(50)      | NULL                      | User's last name                    |
| email_verified     | boolean          | DEFAULT false             | Email verification status           |
| user_type          | varchar(20)      | NOT NULL                  | User type (User/Moderator/Admin/SystemAdmin) |
| verification_token | varchar(100)     | NULL                      | Email verification token            |
| token_expiry       | timestamp        | NULL                      | Verification token expiry time      |
| deletion_requested | boolean          | DEFAULT false             | User has requested account deletion |
| active             | boolean          | DEFAULT true              | Account active status               |
| created_at         | timestamp        | NOT NULL                  | Creation timestamp                  |
| updated_at         | timestamp        | NOT NULL                  | Last update timestamp               |

### Roles

| Column            | Type             | Constraints                 | Description                         |
|-------------------|------------------|----------------------------|-------------------------------------|
| id                | uuid             | PRIMARY KEY                | Role identifier                     |
| name              | varchar(50)      | UNIQUE, NOT NULL          | Role name                           |
| description       | text             | NULL                      | Role description                    |
| created_at        | timestamp        | NOT NULL                  | Creation timestamp                  |
| updated_at        | timestamp        | NOT NULL                  | Last update timestamp               |

### Permissions

| Column            | Type             | Constraints                 | Description                         |
|-------------------|------------------|----------------------------|-------------------------------------|
| id                | uuid             | PRIMARY KEY                | Permission identifier               |
| name              | varchar(100)     | UNIQUE, NOT NULL          | Permission name                     |
| description       | text             | NULL                      | Permission description              |
| resource          | varchar(100)     | NOT NULL                  | Resource being accessed             |
| action            | varchar(50)      | NOT NULL                  | Action on resource (read, write)    |
| created_at        | timestamp        | NOT NULL                  | Creation timestamp                  |
| updated_at        | timestamp        | NOT NULL                  | Last update timestamp               |

### UserRoles

| Column            | Type             | Constraints                     | Description                     |
|-------------------|------------------|--------------------------------|---------------------------------|
| user_id           | uuid             | FK -> users.id, NOT NULL       | Reference to user               |
| role_id           | uuid             | FK -> roles.id, NOT NULL       | Reference to role               |
| assigned_by       | uuid             | FK -> users.id, NOT NULL       | Admin who assigned the role     |
| created_at        | timestamp        | NOT NULL                      | Assignment timestamp            |
| PRIMARY KEY       | (user_id, role_id) |                              | Composite primary key           |

### RolePermissions

| Column            | Type             | Constraints                        | Description                     |
|-------------------|------------------|-----------------------------------|---------------------------------|
| role_id           | uuid             | FK -> roles.id, NOT NULL          | Reference to role               |
| permission_id     | uuid             | FK -> permissions.id, NOT NULL    | Reference to permission         |
| created_at        | timestamp        | NOT NULL                         | Assignment timestamp            |
| PRIMARY KEY       | (role_id, permission_id) |                           | Composite primary key           |

## API Endpoints

| Method | Endpoint                              | Description                  | Access Level       |
|--------|---------------------------------------|------------------------------|-------------------|
| POST   | /api/v1/auth/login                    | Authenticate user            | Public            |
| POST   | /api/v1/auth/logout                   | Invalidate session           | Authenticated     |
| POST   | /api/v1/auth/register                 | Register a user              | Public            |
| GET    | /api/v1/auth/verify                   | Verify email with token      | Public            |
| POST   | /api/v1/auth/resend-verification      | Resend verification email    | Public            |
| POST   | /api/v1/auth/password-reset           | Reset password               | Public            |
| GET    | /api/v1/users                         | List all users               | Admin+            |
| GET    | /api/v1/users/{user_id}               | Get user details             | User (self only)/Moderator+ |
| PUT    | /api/v1/users/{user_id}               | Update user                  | User (self only)/Admin+ |
| POST   | /api/v1/users/{user_id}/request-deletion | Request account deletion  | User (self only)  |
| DELETE | /api/v1/users/{user_id}               | Delete user                  | Moderator+        |
| POST   | /api/v1/users/{user_id}/role          | Change user role             | Admin+            |
| POST   | /api/v1/users/{user_id}/promote/admin | Promote to admin             | System Admin      |
| POST   | /api/v1/users/{user_id}/promote/moderator | Promote to moderator     | Admin+            |
| POST   | /api/v1/users/{user_id}/demote        | Demote user role             | Admin+            |
| GET    | /api/v1/roles                         | List all roles               | Admin+            |
| GET    | /api/v1/roles/{role_id}               | Get role details             | Admin+            |
| POST   | /api/v1/roles                         | Create role                  | Admin+            |
| PUT    | /api/v1/roles/{role_id}               | Update role                  | Admin+            |
| DELETE | /api/v1/roles/{role_id}               | Delete role                  | Admin+            |
| GET    | /api/v1/permissions                   | List all permissions         | Admin+            |
| GET    | /api/v1/permissions/{permission_id}   | Get permission details       | Admin+            |
| GET    | /api/v1/me                            | Get current user profile     | Authenticated     |
| GET    | /api/v1/me/permissions                | Get current user permissions | Authenticated     |

## Getting Started

### Prerequisites
- Go 1.20+
- Docker
- PostgreSQL 13+
- Git

### Setup Instructions

```bash
# Clone the repository
git clone https://github.com/your-username/affpilot-auth-service.git
cd affpilot-auth-service

# Configure environment variables
cp .env.example .env

# Start the application
go run ./cmd/main.go
```

Or using Docker:

```bash
docker-compose up --build
```

## Environment Variables

| Variable                | Description                                |
|-------------------------|--------------------------------------------|
| APP_ENV                 | Application environment                    |
| DB_HOST                 | PostgreSQL host                            |
| DB_PORT                 | PostgreSQL port                            |
| DB_USER                 | PostgreSQL username                        |
| DB_PASSWORD             | PostgreSQL password                        |
| DB_NAME                 | PostgreSQL database name                   |
| JWT_SECRET              | Secret key for signing JWT tokens          |
| JWT_EXPIRY              | JWT token validity duration (e.g. "24h")   |
| SYSTEM_ADMIN_USERNAME   | Initial system admin username              |
| SYSTEM_ADMIN_PASSWORD   | Initial system admin password              |
| SYSTEM_ADMIN_EMAIL      | Initial system admin email                 |
| PASSWORD_SALT           | Salt for password hashing                  |
| EMAIL_VERIFICATION_URL  | Base URL for email verification links      |
| EMAIL_FROM              | Sender email address for system emails     |
| EMAIL_HOST              | SMTP server host                           |
| EMAIL_PORT              | SMTP server port                           |
| EMAIL_USERNAME          | SMTP server username                       |
| EMAIL_PASSWORD          | SMTP server password                       |
| EMAIL_SECURE            | Use TLS for SMTP (true/false)              |
| VERIFICATION_TOKEN_TTL  | Verification token lifetime in minutes     |
| LOG_LEVEL               | Logging level (debug, info, warn, error)   |
| SERVER_PORT             | HTTP server port                           |

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Assignment Deliverables

Candidates should implement:

1. Database migrations for the schema defined above
2. User authentication system with JWT-based sessions
3. Email verification workflow with time-limited tokens (5-minute expiry)
4. Role-based authorization system with four user types:
   - System Admin
   - Admin
   - Moderator
   - User
5. System administrator initialization from environment variables
6. Account deletion request workflow
7. API endpoints with proper access control
8. Unit and integration tests
9. Documentation for any additional design decisions

## User Types and Permissions

### System Admin
- Created when the database is initialized (from environment variables)
- Has complete access to all system functions
- Cannot be deleted from the system
- Can promote users to any role
- Can demote any user except other System Admins
- Can manage all roles and permissions

### Admin
- Has full administrative access
- Cannot remove System Admins or other Admins
- Cannot promote users to System Admin role
- Can promote users to Moderator or Admin
- Can manage all roles and permissions

### Moderator
- Can create, read, update, and delete data
- Cannot delete user accounts directly
- Can process deletion requests from users
- Can manage content but not system configuration

### User
- Can only create, read, update, and delete their own data
- Cannot access other users' data
- Can request account deletion (to be processed by Moderator+)

## Email Verification Process

1. When a user registers, their account is created with `email_verified = false`
2. A verification token is generated and stored in the user record
3. The token expiry is set to 5 minutes from creation time
4. An email is sent to the user with a verification link containing the token
5. When the user clicks the link within 5 minutes:
   - The token is validated
   - The user's email is marked as verified
   - The token is cleared from the database
6. If the token has expired:
   - The user is prompted to request a new verification email
   - A new token is generated with a new 5-minute expiry

## Account Deletion Workflow

1. User requests account deletion via API
2. User record is flagged with `deletion_requested = true`
3. Moderator, Admin, or System Admin reviews the request
4. Upon approval, the user account is deleted from the system

## License

This project is licensed under the MIT License.

## Acknowledgements

- [Golang](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [JWT](https://jwt.io/)
EOF

# Create .env.example file
cat > .env.example << 'EOF'
# Application Environment
APP_ENV=development
SERVER_PORT=8080
LOG_LEVEL=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=affpilot_auth

# JWT Configuration
JWT_SECRET=your-secret-key-here
JWT_EXPIRY=24h

# System Admin Configuration
SYSTEM_ADMIN_USERNAME=admin
SYSTEM_ADMIN_PASSWORD=adminpassword
SYSTEM_ADMIN_EMAIL=admin@example.com

# Security
PASSWORD_SALT=your-password-salt-here

# Email Configuration
EMAIL_VERIFICATION_URL=http://localhost:8080/api/v1/auth/verify
EMAIL_FROM=no-reply@example.com
EMAIL_HOST=smtp.example.com
EMAIL_PORT=587
EMAIL_USERNAME=smtp-user
EMAIL_PASSWORD=smtp-password
EMAIL_SECURE=true
VERIFICATION_TOKEN_TTL=5
EOF

# Create .gitignore
cat > .gitignore << 'EOF'
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# Environment files
.env

# IDE files
.idea/
.vscode/
*.swp
*.swo

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Logs
logs/
*.log
EOF

# Create go.mod file
cat > go.mod << 'EOF'
module github.com/your-username/affpilot-auth-service

go 1.20

require (
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.10.0
)
EOF

# Create docker-compose.yml
cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      - APP_ENV=development
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=affpilot_auth
      - JWT_SECRET=your-secret-key-here
      - JWT_EXPIRY=24h
      - SYSTEM_ADMIN_USERNAME=admin
      - SYSTEM_ADMIN_PASSWORD=adminpassword
      - SYSTEM_ADMIN_EMAIL=admin@example.com
      - PASSWORD_SALT=your-password-salt-here
      - EMAIL_VERIFICATION_URL=http://localhost:8080/api/v1/auth/verify
      - EMAIL_FROM=no-reply@example.com
      - EMAIL_HOST=smtp.example.com
      - EMAIL_PORT=587
      - EMAIL_USERNAME=smtp-user
      - EMAIL_PASSWORD=smtp-password
      - EMAIL_SECURE=true
      - VERIFICATION_TOKEN_TTL=5
      - LOG_LEVEL=debug
      - SERVER_PORT=8080
  
  postgres:
    image: postgres:14-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=affpilot_auth
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
EOF

# Create a simple Dockerfile
cat > Dockerfile << 'EOF'
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

EXPOSE 8080

CMD ["./main"]
EOF

# Create a simple database migration file
mkdir -p migrations/000001_init_schema
cat > migrations/000001_init_schema/up.sql << 'EOF'
-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email_verified BOOLEAN DEFAULT FALSE,
    user_type VARCHAR(20) NOT NULL,
    verification_token VARCHAR(100),
    token_expiry TIMESTAMP,
    deletion_requested BOOLEAN DEFAULT FALSE,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- UserRoles junction table
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

-- RolePermissions junction table
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

-- Insert default roles
INSERT INTO roles (name, description) VALUES
    ('system_admin', 'Full system access with ability to manage all aspects of the system'),
    ('admin', 'Administrative access with limitations on managing other admins'),
    ('moderator', 'Can manage content but cannot delete users directly'),
    ('user', 'Regular user with access only to their own data');

-- Insert default permissions
INSERT INTO permissions (name, resource, action, description) VALUES
    ('user:read:all', 'user', 'read:all', 'Read all user data'),
    ('user:create:all', 'user', 'create:all', 'Create any user'),
    ('user:update:all', 'user', 'update:all', 'Update any user'),
    ('user:delete:all', 'user', 'delete:all', 'Delete any user'),
    ('user:read:self', 'user', 'read:self', 'Read own user data'),
    ('user:update:self', 'user', 'update:self', 'Update own user data'),
    ('user:delete:self', 'user', 'delete:self', 'Request own account deletion'),
    ('role:read', 'role', 'read', 'Read roles'),
    ('role:create', 'role', 'create', 'Create roles'),
    ('role:update', 'role', 'update', 'Update roles'),
    ('role:delete', 'role', 'delete', 'Delete roles'),
    ('permission:read', 'permission', 'read', 'Read permissions'),
    ('user:promote:admin', 'user', 'promote:admin', 'Promote user to admin'),
    ('user:promote:moderator', 'user', 'promote:moderator', 'Promote user to moderator'),
    ('user:demote', 'user', 'demote', 'Demote user role');

-- Assign permissions to roles
-- System Admin permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'system_admin'), 
    id 
FROM permissions;

-- Admin permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'admin'), 
    id 
FROM permissions
WHERE name NOT IN ('user:promote:admin');

-- Moderator permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'moderator'), 
    id 
FROM permissions
WHERE name IN (
    'user:read:all', 
    'user:read:self', 
    'user:update:self', 
    'user:delete:self'
);

-- User permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'user'), 
    id 
FROM permissions
WHERE name IN (
    'user:read:self', 
    'user:update:self', 
    'user:delete:self'
);
EOF

cat > migrations/000001_init_schema/down.sql << 'EOF'
-- Drop tables in reverse order
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
EOF

# Create a basic config loader in Go
mkdir -p internal/config
cat > internal/config/config.go << 'EOF'
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Admin    AdminConfig
	Email    EmailConfig
	Server   ServerConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Environment string
	LogLevel    string
}

// DatabaseConfig holds database connection information
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// AdminConfig holds system admin information
type AdminConfig struct {
	Username string
	Password string
	Email    string
}

// EmailConfig holds email configuration
type EmailConfig struct {
	VerificationURL  string
	From             string
	Host             string
	Port             int
	Username         string
	Password         string
	Secure           bool
	VerificationTTL  int
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Parse DB port
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	
	// Parse JWT expiry
	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	
	// Parse email port
	emailPort, _ := strconv.Atoi(getEnv("EMAIL_PORT", "587"))
	
	// Parse email secure
	emailSecure, _ := strconv.ParseBool(getEnv("EMAIL_SECURE", "true"))
	
	// Parse verification token TTL
	verificationTTL, _ := strconv.Atoi(getEnv("VERIFICATION_TOKEN_TTL", "5"))
	
	// Parse server port
	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))

	return &Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "affpilot_auth"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key-here"),
			Expiry: jwtExpiry,
		},
		Admin: AdminConfig{
			Username: getEnv("SYSTEM_ADMIN_USERNAME", "admin"),
			Password: getEnv("SYSTEM_ADMIN_PASSWORD", "adminpassword"),
			Email:    getEnv("SYSTEM_ADMIN_EMAIL", "admin@example.com"),
		},
		Email: EmailConfig{
			VerificationURL: getEnv("EMAIL_VERIFICATION_URL", "http://localhost:8080/api/v1/auth/verify"),
			From:            getEnv("EMAIL_FROM", "no-reply@example.com"),
			Host:            getEnv("EMAIL_HOST", "smtp.example.com"),
			Port:            emailPort,
			Username:        getEnv("EMAIL_USERNAME", ""),
			Password:        getEnv("EMAIL_PASSWORD", ""),
			Secure:          emailSecure,
			VerificationTTL: verificationTTL,
		},
		Server: ServerConfig{
			Port: serverPort,
		},
	}, nil
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
EOF

echo "Project structure has been created successfully!"