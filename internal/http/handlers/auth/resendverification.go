package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"log"

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

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer database.Close(db)

	// Generate a new verification token
	newToken := uuid.NewString()

	// Update the verification token in the database
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

	// Send the verification email
	if err := sendVerificationEmail(reqBody.Email, newToken); err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		log.Println("Failed to send verification email to:", reqBody.Email, "Error:", err)
		return
	}

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
