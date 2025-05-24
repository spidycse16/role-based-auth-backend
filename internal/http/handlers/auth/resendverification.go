package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
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

	db:=database.Connect()

	// Step 1: Check if the user exists and if email is already verified
	var emailVerified bool
	err := db.QueryRow("SELECT email_verified FROM users WHERE email = $1", reqBody.Email).Scan(&emailVerified)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"status":  strconv.Itoa(http.StatusOK),
			"message": "Email already verified",
			"data":    nil,
		}
		json.NewEncoder(w).Encode(response)
		log.Println("Email already verified for:", reqBody.Email)
		return
	}

	// Step 2: Generate a new verification token
	newToken := uuid.NewString()

	// Step 3: Update the verification token in the database
	_, err = db.Exec(
		"UPDATE users SET verification_token = $1, updated_at = NOW() WHERE email = $2 AND email_verified = false",
		newToken,
		reqBody.Email,
	)
	if err != nil {
		http.Error(w, "Failed to update verification token", http.StatusInternalServerError)
		log.Println("Failed to update verification token for:", reqBody.Email, "Error:", err)
		return
	}

	// Step 4: Send the verification email
	if err := sendVerificationEmail(reqBody.Email, newToken); err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		log.Println("Failed to send verification email to:", reqBody.Email, "Error:", err)
		return
	}

	// Step 5: Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Verification email sent successfully",
		"data":    nil,
	}
	json.NewEncoder(w).Encode(response)

	log.Println("Verification email resent successfully to:", reqBody.Email)
}
