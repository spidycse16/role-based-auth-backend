package handlers


import (
	"encoding/json"
	"database/sql"
	"net/http"
	"strings"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/models"
	"github.com/google/uuid"
	"regexp"
)

type RegisterRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !isValidEmail(req.Email) {
		http.Error(w, "Invalid Email foramt", http.StatusBadRequest)
		return
	}
	

	// Trim spaces
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Password = strings.TrimSpace(req.Password)

	// Validate required fields
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Generate verification token and expiry
	verificationToken := uuid.NewString()
	tokenExpiry := time.Now().Add(24 * time.Hour)

	// Create the user model
	user := models.User{
		Username:          req.Username,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Email:             req.Email,
		PasswordHash:      string(hashedPassword),
		EmailVerified:     false,
		UserType:          "User",
		Active:            true,
		VerificationToken: verificationToken,
		TokenExpiry:       &tokenExpiry, // Use the address here
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Insert the user into the database
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer database.Close(db)

	// Check if the email already exists
	if exists,err := isEmailExists(db, user.Email); err != nil {
		http.Error(w, "Failed to check email existence", http.StatusInternalServerError)
		return
	} else if exists {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}
	query := `
	INSERT INTO users (username, email, password_hash, first_name, last_name, email_verified, user_type, active, verification_token, token_expiry, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id
	`

	err = db.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.EmailVerified, user.UserType, user.Active, user.VerificationToken, user.TokenExpiry, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Send verification email
	if err := sendVerificationEmail(req.Email, verificationToken); err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		return
	}

	// Return the created user (without the password)
	user.PasswordHash = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// isEmailExists checks if an email is already registered.
func isEmailExists(db *sql.DB, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

// isValidEmail validates the email format.
func isValidEmail(email string) bool {
	// Basic email regex, adjust if needed
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
