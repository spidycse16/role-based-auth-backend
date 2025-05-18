package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/sagorsarker04/Developer-Assignment/internal/config"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the JSON response for a successful login.
type LoginResponse struct {
	Token string `json:"token"`
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
	var emailVerified bool
	query := "SELECT id, username, user_type, password_hash, email_verified FROM users WHERE email = $1"
	err = db.QueryRow(query, req.Email).Scan(&userID, &username, &userType, &storedHash, &emailVerified)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	// Check if the email is verified
	if !emailVerified {
		http.Error(w, "Email not verified", http.StatusForbidden)
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
	expiresAt := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Expires:  expiresAt,
		Secure:   false, // Set to true in production with HTTPS
		Path:     "/",
	}
	// Set cookie in response
	http.SetCookie(w, cookie)
	w.Write([]byte("Login successful! Cookie set."))

	// Return the token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
