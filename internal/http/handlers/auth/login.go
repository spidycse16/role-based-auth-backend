package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/sagorsarker04/Developer-Assignment/internal/config"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Type     string `json:"type"`
}

// LoginResponse represents the JSON response for a successful login.
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// http.Error(w, "Invalid request payload", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request Payload")
		return
	}

	// Trim spaces
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		// http.Error(w, "Email and password are required", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Load the config
	cfg:=config.GetConfig()

	// Connect to the database
	db:=database.Connect()

	// Fetch user from the database
	var storedHash, userID, username, userType string
	var emailVerified bool
	query := "SELECT id, username, user_type, password_hash, email_verified FROM users WHERE email = $1"
	err := db.QueryRow(query, req.Email).Scan(&userID, &username, &userType, &storedHash, &emailVerified)
	if err == sql.ErrNoRows {
		// http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	} else if err != nil {
		// http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	// Check if the email is verified
	if !emailVerified {
		// http.Error(w, "Email not verified", http.StatusForbidden)
		utils.ErrorResponse(w, http.StatusForbidden, "Email not verified")
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
		// http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := generateJWT(userID, username, userType, cfg.JWT.Secret, cfg.JWT.Expiry)
	if err != nil {
		// http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
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
	
	// Set header and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	// Prepare your JSON response manually
	response := struct {
		Status  string        `json:"status"`
		Message string        `json:"message"`
		Data    LoginResponse `json:"data"`
	}{
		Status:  "202",
		Message: "Login successful",
		Data: LoginResponse{
			Token: token,
			User: UserInfo{
				ID:       userID,
				Username: username,
				Email:    req.Email,
				Type:     userType,
			},
		},
	}

	// Encode response to JSON and write it to response body
	json.NewEncoder(w).Encode(response)
}