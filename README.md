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
