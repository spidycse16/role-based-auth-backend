package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sagorsarker04/Developer-Assignment/internal/config"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
)

// ResendVerificationEmail handles resending the verification email to unverified users.
func ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email string `json:"email"`
	}

	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Invalid request body:", err)
		return
	}

	db := database.Connect()

	// Step 1: Check if user exists and get user_id, email_verified
	var userID uuid.UUID
	var emailVerified bool
	err := db.QueryRow(
		"SELECT id, email_verified FROM users WHERE email = $1",
		reqBody.Email,
	).Scan(&userID, &emailVerified)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User with this email does not exist", http.StatusNotFound)
			log.Println("No user found with email:", reqBody.Email)
			return
		}
		http.Error(w, "Database query error", http.StatusInternalServerError)
		log.Println("Query error for email:", reqBody.Email, "Error:", err)
		return
	}

	if emailVerified {
		response := map[string]interface{}{
			"status":  strconv.Itoa(http.StatusOK),
			"message": "Email already verified",
			"data":    nil,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		log.Println("Email already verified for:", reqBody.Email)
		return
	}

	// Step 2: Generate a new verification token
	cfg := config.GetConfig()
	secretKey := cfg.JWT.Secret
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   reqBody.Email,
		"exp":     time.Now().Add(cfg.Email.VerificationTTL).Unix(),
		"iat":     time.Now().Unix(),
		"purpose": "email_verification",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	verificationToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		http.Error(w, "Failed to generate verification token", http.StatusInternalServerError)
		log.Println("Token generation failed:", err)
		return
	}

	// Step 3: Update the verification token in database (optional field)
	_, err = db.Exec(
		"UPDATE users SET updated_at = NOW() WHERE email = $1 AND email_verified = false",
		reqBody.Email,
	)
	if err != nil {
		http.Error(w, "Failed to update verification token", http.StatusInternalServerError)
		log.Println("Failed to update verification token for:", reqBody.Email, "Error:", err)
		return
	}

	// Step 4: Send the verification email
	if err := sendVerificationEmail(reqBody.Email, verificationToken); err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		log.Println("Failed to send verification email to:", reqBody.Email, "Error:", err)
		return
	}

	// Step 5: Return success response
	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Verification email sent successfully",
		"data":    nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	log.Println("Verification email resent successfully to:", reqBody.Email)
}
