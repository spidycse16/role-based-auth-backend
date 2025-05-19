# Developer Assignment Documentation

## Role Based Auth System

This is a robust, backend-only authentication and user management system built with Golang, designed to provide secure, scalable API-driven solutions. It leverages JWT-based authentication, role-based access control (RBAC), and email verification, making it ideal for startups, SaaS platforms, internal dashboards, community sites, educational tools, and e-commerce applications. The project is structured for maintainability and deployed easily with Docker.

### Tech Stack

- **Backend**: Golang (net/http, Gorilla Mux)(Need for local installation)
- **Database**: PostgreSQL (Need for local installation)
- **Deployment**: Docker, Docker Compose (Docker Based Setup)

## Installation and Setup

To get started with the project, follow these detailed steps. For the easiest and most reliable installation, we strongly recommend using Docker, which leverages the provided `Dockerfile` and `docker-compose.yml` to handle dependencies and setup.

You can just Clone the project copy the .env.example and run the project by docker using docker compose up --build.

### 1. Prerequisites

Ensure the following software is installed on your system:

- **Go**: v1.20 or higher (required for local development, optional with Docker)
- **PostgreSQL**: v15 or higher (optional if using Docker)
- **Docker**: v24 or higher (recommended for easy installation)
- **Docker Compose**: v2.20 or higher (included with Docker)
- **Git**: v2.40 or higher

**Note**: Installing Docker simplifies the process by containerizing the application and database, eliminating the need to manually install Go or PostgreSQL on your host machine. Download and install Docker from docker.com if you haven’t already.

### 2. Clone the Repository

Clone the project repository and navigate into the project directory:

```bash
git clone https://github.com/sagorsarker04/Developer-Assignment/tree/Development
cd Developer-Assignment
```

### 3. Configure Environment Variables

Copy the `.env.example` file to `.env` and modify the following sections with your specific details:

```bash
cp .env.example .env
```

#### Database Configuration

Update the database settings to match your PostgreSQL instance or the Dockerized PostgreSQL service:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=your_database_name
```

- If using Docker (recommended), leave `DB_HOST=postgres` as configured in `docker-compose.yml`, and ensure `DB_PASSWORD` and `DB_NAME` match the `postgres` service settings (e.g., `123` and `affpilot_auth`).

#### Email Configuration

Configure the email service for verification and password reset emails. Replace with your SMTP provider details:

```
SMTP_HOST=smtp.yourprovider.com
SMTP_PORT=587
SMTP_USER=your_email@example.com
SMTP_PASSWORD=your_email_password
EMAIL_FROM=your_email@example.com
EMAIL_VERIFICATION_URL=http://localhost:8080/api/v1/auth/verify
```

- Use a service like Gmail (e.g., `smtp.gmail.com`) and generate an App Password if two-factor authentication is enabled.
- Ensure `EMAIL_VERIFICATION_URL` points to your running application’s URL.

#### Other Configurations

- Leave `JWT_SECRET` as a secure, random string (e.g., `your_secret_key`).
- Set `JWT_EXPIRY` to your preferred token duration (e.g., `24h`).
- Set `APP_PORT=8080` unless you need a different port.

### 4. Run Migrations (First Time Only)

The project uses Docker for deployment. To initialize the database with migrations:

1. Open the `Dockerfile` and uncomment the migration command:

   ```dockerfile
   CMD ["go", "run", "cmd/migration/main.go"]
   ```
2. Comment out the application command:

   ```dockerfile
   CMD ["go", "run", "cmd/main.go"]
   ```
3. Build and run the Docker containers:

   ```bash
   docker-compose up --build
   ```
4. After migrations complete (check logs for confirmation), You have successfully created all the tables and seeded basic infos in the tables like table migration, default roles, user roles and the default permissons for each roles. And we are ready to start the project and use its feature
now, stop the containers:

   ```bash
   docker-compose down
   ```
5. Revert the `Dockerfile` changes: comment the migration line and uncomment the CMD ["go", "run", "cmd/main.go"] line.

### 5. Run the Application

Start the application using Docker:

```bash
docker-compose up --build
```

The backend will be accessible at `http://localhost:8080`. Use a tool like Postman to interact with the APIs (see Postman Collection).

### 6. Set Up System Admin

Since this is a backend-only project, the System Admin role is not set automatically. Follow these steps:

1. Register a new user using the `/auth/register` API endpoint.
2. Access the PostgreSQL database (e.g., via `docker-compose exec postgres psql -U postgres`).
3. Update the `user_type` field of the registered user to `system_admin` in the `users` table:

   ```sql
   UPDATE users SET user_type = 'system_admin' WHERE id = your_user_id;
   ```

## Project Overview

The Developer Assignment project is a sophisticated backend solution crafted to address the needs of modern web applications requiring secure user management and authentication. Built entirely with Golang, it utilizes the net/http package and Gorilla Mux router to create a modular and efficient API framework. The system emphasizes security through features like JWT-based authentication stored in HTTP-only cookies, ensuring protection against common vulnerabilities such as XSS attacks. Role-Based Access Control (RBAC) is a cornerstone of the project, supporting a multi-tier role hierarchy including System Admin, Admin, Moderator, and User, each with finely tuned permissions to prevent unauthorized access.

Key security practices include password hashing with bcrypt, token expiry to limit attack windows, and email verification to ensure only validated users gain full access. The project supports advanced account management features such as user registration, login, profile updates, role promotions/demotions, and secure password resets. Its modular design separates concerns into distinct handlers, middleware, and database logic, promoting clean code and easy maintenance.

The tech stack includes PostgreSQL for robust database management, custom middleware for RBAC enforcement, and Docker with Docker Compose for seamless deployment. This backend-only system is designed for API interaction, with a provided Postman collection to facilitate testing and development. Use cases span a wide range, from multi-tenant SaaS platforms to internal company tools and community-driven applications, offering flexibility and scalability for diverse needs.

### Key Features

- **JWT-Based Authentication**: Implements stateless, secure sessions using JSON Web Tokens stored in HTTP-only cookies with expiration settings.
- **Role-Based Access Control (RBAC)**: Offers a multi-tier role system (System Admin, Admin, Moderator, User) with strict hierarchy enforcement and granular permission checks.
- **Email Verification**: Mandates email verification with unique tokens, including a resend feature for unverified users.
- **Secure Password Handling**: Utilizes bcrypt for strong password hashing and storage.
- **Account Management**: Provides APIs for registration, login, updates, role management, and password resets with token-based recovery.
- **Demotion and Promotion**: Includes APIs to adjust user roles while respecting hierarchical constraints.
- **Modular Design**: Features a clean, maintainable code structure with separated concerns for handlers, middleware, and database interactions.

## API Endpoints

This is a backend-only project, designed for API interaction. All endpoints are prefixed with the base URL `http://localhost:8080/api/v1/`. For example, the `/users` endpoint is accessible at `http://localhost:8080/api/v1/users`. Use the provided Postman collection or manually test the endpoints. Responses follow a standardized format:

### Example Response Format

```json
{
  "status": "200",
  "message": "Users retrieved successfully",
  "data": [...]
}
```

### Authentication

| Endpoint | Method | Description | Authentication Required |
| --- | --- | --- | --- |
| `http://localhost:8080/api/v1/auth/login` | POST | Log in a user | No |
| `http://localhost:8080/api/v1/auth/logout` | POST | Log out a user | Yes |
| `http://localhost:8080/api/v1/auth/register` | POST | Register a new user | No |
| `http://localhost:8080/api/v1/auth/verify/{token}` | GET | Verify email with token | No |
| `http://localhost:8080/api/v1/auth/resend-verification` | POST | Resend email verification link | Yes |
| `http://localhost:8080/api/v1/auth/password-reset-request` | POST | Request a password reset | No |
| `http://localhost:8080/api/v1/auth/password-reset-confirm` | POST | Confirm password reset | No |

### Permissions

| Endpoint | Method | Description | Authentication Required | Role Requirement |
| --- | --- | --- | --- | --- |
| `http://localhost:8080/api/v1/permissions/create` | POST | Create a new permission | Yes | Admin+ |
| `http://localhost:8080/api/v1/permissions` | GET | List all permissions | Yes | Admin+ |
| `http://localhost:8080/api/v1/permissions/{permission_id}` | GET | Get permission details | Yes | Admin+ |

### Current User

| Endpoint | Method | Description | Authentication Required |
| --- | --- | --- | --- |
| `http://localhost:8080/api/v1/me` | GET | Get current user profile | Yes |
| `http://localhost:8080/api/v1/me/permissions` | GET | Get current user permissions | Yes |

### Roles

| Endpoint | Method | Description | Authentication Required | Role Requirement |
| --- | --- | --- | --- | --- |
| `http://localhost:8080/api/v1/roles` | GET | List all roles | Yes | `role:read`, `admin:read`, or `system_admin:read` |
| `http://localhost:8080/api/v1/roles/{role_id}` | GET | Get role details | Yes | None |
| `http://localhost:8080/api/v1/roles/create` | POST | Create a new role | Yes | Admin+ |
| `http://localhost:8080/api/v1/roles/{role_id}` | PUT | Update a role | Yes | Admin+ |
| `http://localhost:8080/api/v1/roles/{role_id}` | DELETE | Delete a role | Yes | Admin+ |
| `http://localhost:8080/api/v1/roles/{user_id}/role` | POST | Change a user’s role | Yes | Admin+ |
| `http://localhost:8080/api/v1/roles/{user_id}/promote/admin` | POST | Promote user to Admin | Yes | `user:promote:admin` |
| `http://localhost:8080/api/v1/roles/{user_id}/promote/moderator` | POST | Promote user to Moderator | Yes | `user:promote:moderator` |
| `http://localhost:8080/api/v1/roles/{user_id}/demote` | POST | Demote a user | Yes | `user:demote` |

### Users

| Endpoint | Method | Description | Authentication Required | Role Requirement |
| --- | --- | --- | --- | --- |
| `http://localhost:8080/api/v1/users` | GET | List all users | Yes | `user:read:all` |
| `http://localhost:8080/api/v1/users/{user_id}` | GET | Get user details | Yes | `user:read:self` |
| `http://localhost:8080/api/v1/users/{user_id}` | PUT | Update user details | Yes | `user:update:self` or `user:update:all` |
| `http://localhost:8080/api/v1/users/{user_id}` | POST | Request user deletion (soft delete) | Yes | `user:delete:self` |
| `http://localhost:8080/api/v1/users/{user_id}` | DELETE | Permanently delete a user | Yes | `user:delete:all` |

## Postman Collection

For easy API testing, use the provided Postman collection:\
https://affpilot-9730.postman.co/workspace/My-Workspace\~9b9ff17b-ea06-458a-8ef8-7c9cfa53edc1/collection/44831755-b27c1d36-8339-4670-9a9a-3e0ab57045fb?action=share&creator=44831755

Import this collection into Postman to test the APIs with pre-configured requests.

## Security Best Practices

- **Password Hashing**: Passwords are hashed using bcrypt.
- **Token Expiry**: JWTs have expiration times to reduce attack windows.
- **HTTP-Only Cookies**: JWTs are stored in HTTP-only cookies to prevent XSS attacks.
- **Email Verification**: Unverified accounts have restricted access.
- **Role Hierarchy**: Enforces strict role hierarchies to prevent privilege escalation.

## Project Structure

- `cmd/`: Entry points for the application and migrations.
- `migration/`: Database migration scripts.
- `internal/`: Core application logic.
  - `config/`: Configuration management.
  - `database/`: Database setup and connections.
  - `handlers/`: API route handlers.
    - `auth/`: Authentication handlers.
    - `permission/`: Permission handlers.
    - `role/`: Role management handlers.
    - `user/`: User management handlers.
  - `middleware/`: Custom middleware for authentication and permissions.
  - `routes/`: API route definitions.
- `models/`: Data models for the application.
- `services/`: Business logic and services.
- `migrations/`: Migration files for database schema changes.
- `scripts/`: Utility scripts.
- `static/`: Static assets (if any).
- `test/`: Test files.
- `.env`: Environment variables.
- `.env.example`: Example environment variables.
- `.gitignore`: Git ignore file.
- `docker-compose.yml`: Docker Compose configuration.
- `Dockerfile`: Docker configuration for the application.
- `go.mod`, `go.sum`: Go module dependencies.
- `README.md`: Project readme.
- `setup.sh`: Setup script (if any).

## Use Cases

- **Startups and SaaS Solutions**: Ideal for multi-tenant platforms with tiered access control.
- **Internal Dashboards**: Suitable for tools with varying employee access levels.
- **Community Platforms**: Supports role-based moderation for forums and social networks.
- **Educational Platforms**: Manages instructor, student, and admin roles.
- **E-commerce Marketplaces**: Handles vendors, buyers, and admins with distinct permissions.