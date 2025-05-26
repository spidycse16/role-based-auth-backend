package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sagorsarker04/Developer-Assignment/internal/config"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/models"
	"golang.org/x/crypto/bcrypt"
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

	if len(req.Password) < 8 || len(req.Password) > 20 {
		http.Error(w, "Password should be 8 to 20 characters long", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Generate verification token and expiry
	// verificationToken := uuid.NewString()
	// tokenExpiry := time.Now().Add(24 * time.Hour)

	// Create the user model
	user := models.User{
		Username:      req.Username,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         req.Email,
		PasswordHash:  string(hashedPassword),
		EmailVerified: false,
		UserType:      "user",
		Active:        true,
		ResetToken:    "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Insert the user into the database
	db := database.Connect()

	// Check if the email already exists
	if exists, err := isEmailExists(db, user.Email); err != nil {
		http.Error(w, "Failed to check email existence", http.StatusInternalServerError)
		return
	} else if exists {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}
	query := `
	INSERT INTO users (username, email, password_hash, first_name, last_name, email_verified, user_type, active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
	`

	err = db.QueryRow(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.EmailVerified,
		user.UserType,
		user.Active,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	cfg := config.GetConfig()
	secretKey := cfg.JWT.Secret
	// Create verification token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(cfg.Email.VerificationTTL).Unix(),
		"iat":     time.Now().Unix(),
		"purpose": "email_verification",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	verificationToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		http.Error(w, "Failed to generate verification token", http.StatusInternalServerError)
		return
	}
	// Send verification email
	if err := sendVerificationEmail(req.Email, verificationToken); err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		return
	}

	user.PasswordHash = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "User created successfully",
		"data":    user,
	}

	json.NewEncoder(w).Encode(response)

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


