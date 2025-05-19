package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strconv"

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

	// Assign default "User" role
	defaultRoleID := "7f6fa0ea-b43d-4443-8be5-d2ce1ff663e9"
	systemAdminID := "66557562-d184-49ad-83ac-9dc3841ed365"

	// Check if the role is already assigned to avoid duplicate entries
	var roleExists bool
	roleCheckQuery := `
		SELECT EXISTS(
			SELECT 1 FROM user_roles 
			WHERE user_id = $1 AND role_id = $2
		)
	`
	err = db.QueryRow(roleCheckQuery, userID, defaultRoleID).Scan(&roleExists)
	if err != nil {
		http.Error(w, "Failed to check user role", http.StatusInternalServerError)
		log.Println("Failed to check user role for user:", userID, "Error:", err)
		return
	}

	// If role not assigned, assign it
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Email verified successfully",
		"data":    nil, // no extra data to send
	}

	json.NewEncoder(w).Encode(response)

	log.Println("Email verified and role assigned successfully for user:", userID)

}

// func printToken(){

// }
