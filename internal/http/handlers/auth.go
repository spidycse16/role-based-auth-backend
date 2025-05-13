package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"
	"regexp" 
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/sagorsarker04/Developer-Assignment/internal/config"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/models"
	
	"log"
	"github.com/gorilla/mux"
)


// RegisterRequest represents the expected JSON payload for user registration.
type RegisterRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the JSON response for a successful login.
type LoginResponse struct {
	Token string `json:"token"`
}
// RegisterUser handles user registration.
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
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
		TokenExpiry:       &tokenExpiry,  // Use the address here
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


// sendVerificationEmail sends a verification email to the user.
func sendVerificationEmail(email, token string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	verificationURL := fmt.Sprintf("%s/%s", cfg.Email.VerificationURL, token)

	// Set up SMTP
	auth := smtp.PlainAuth("", cfg.Email.Username, cfg.Email.Password, cfg.Email.Host)
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: Verify your email\n\nClick the link to verify your email: %s", cfg.Email.From, email, verificationURL)

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", cfg.Email.Host, cfg.Email.Port),
		auth,
		cfg.Email.From,
		[]string{email},
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
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




func LoginUser(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Trim spaces
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Load the config
	cfg, err := config.LoadConfig()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer database.Close(db)

	// Fetch user from the database
	var storedHash, userID, username, userType string
	query := "SELECT id, username, user_type, password_hash FROM users WHERE email = $1"
	err = db.QueryRow(query, req.Email).Scan(&userID, &username, &userType, &storedHash)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateJWT(userID, username, userType, cfg.JWT.Secret, cfg.JWT.Expiry)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

// generateJWT generates a JWT token for the authenticated user.
func generateJWT(userID, username, userType, secret string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"username":  username,
		"user_type": userType,
		"exp":       time.Now().Add(expiry).Unix(),
		"iat":       time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

//returns the global token

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := strings.TrimSpace(vars["token"])
	log.Printf("Verification Token: %s", token)

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer database.Close(db)

	// Update the user's email verification status
	query := `
		UPDATE users 
		SET email_verified = true, updated_at = NOW() 
		WHERE verification_token = $1 AND email_verified = false
		RETURNING id
	`
	var userID string
	err = db.QueryRow(query, token).Scan(&userID)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		log.Println("Invalid or expired token:", token)
		return
	} else if err != nil {
		http.Error(w, "Failed to verify email", http.StatusInternalServerError)
		log.Println("Database error during verification:", err)
		return
	}

	// Clear the verification token to prevent reuse
	_, err = db.Exec("UPDATE users SET verification_token = NULL WHERE id = $1", userID)
	if err != nil {
		http.Error(w, "Failed to clear verification token", http.StatusInternalServerError)
		log.Println("Failed to clear verification token for user:", userID, "Error:", err)
		return
	}

	w.Write([]byte("Email verified successfully"))
}


// func printToken(){

// }

