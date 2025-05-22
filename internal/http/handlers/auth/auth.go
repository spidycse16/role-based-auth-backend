package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"

	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/config"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
)

// sendVerificationEmail sends a verification email to the user.
func sendVerificationEmail(email, token string) error {
	cfg := config.GetConfig()

	verificationURL := fmt.Sprintf("%s/%s", cfg.Email.VerificationURL, token)

	// Set up SMTP
	auth := smtp.PlainAuth("", cfg.Email.Username, cfg.Email.Password, cfg.Email.Host)
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: Verify your email\n\nClick the link to verify your email: %s", cfg.Email.From, email, verificationURL)

	err := smtp.SendMail(
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

	// Step 1: Verify token and get user ID
	var userID string
	query := `
		UPDATE users 
		SET email_verified = true, updated_at = NOW() 
		WHERE verification_token = $1 AND email_verified = false
		RETURNING id
	`
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

	// Step 2: Clear the verification token
	_, err = db.Exec("UPDATE users SET verification_token = NULL WHERE id = $1", userID)
	if err != nil {
		http.Error(w, "Failed to clear verification token", http.StatusInternalServerError)
		log.Println("Failed to clear verification token for user:", userID, "Error:", err)
		return
	}

	// Step 3: Get the ID of the default role ("user")
	var defaultRoleID string
	err = db.QueryRow(`SELECT id FROM roles WHERE name = 'user'`).Scan(&defaultRoleID)
	if err == sql.ErrNoRows {
		http.Error(w, "Default role not found", http.StatusInternalServerError)
		log.Println("Default role 'user' not found")
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch default role", http.StatusInternalServerError)
		log.Println("Failed to get default role ID:", err)
		return
	}

	// Step 4: Get system admin ID (assuming there's one system admin)
	var systemAdminID string
	var systemAdminRoleID string
    err = db.QueryRow(`SELECT id FROM roles WHERE name = 'system_admin'`).Scan(&systemAdminRoleID)

	if err == sql.ErrNoRows {
		http.Error(w, "System admin not found", http.StatusInternalServerError)
		log.Println("System admin not found")
		return
	} else if err != nil {
		http.Error(w, "Failed to find system admin", http.StatusInternalServerError)
		log.Println("Failed to get system admin ID:", err)
		return
	}

	// Step 5: Check if role already assigned
	var roleExists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM user_roles 
			WHERE user_id = $1 AND role_id = $2
		)
	`, userID, defaultRoleID).Scan(&roleExists)
	if err != nil {
		http.Error(w, "Failed to check user role", http.StatusInternalServerError)
		log.Println("Failed to check user role for user:", userID, "Error:", err)
		return
	}

	// Step 6: Assign role if not already assigned
	if !roleExists {
		_, err = db.Exec(`
			INSERT INTO user_roles (user_id, role_id, assigned_by, created_at)
			VALUES ($1, $2, $3, NOW())
		`, userID, defaultRoleID, systemAdminID)
		if err != nil {
			http.Error(w, "Failed to assign default role", http.StatusInternalServerError)
			log.Println("Failed to assign default role for user:", userID, "Error:", err)
			return
		}
	}

	// Step 7: Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Email verified successfully",
		"data":    nil,
	}
	json.NewEncoder(w).Encode(response)
	log.Println("Email verified and role assigned successfully for user:", userID)
}
