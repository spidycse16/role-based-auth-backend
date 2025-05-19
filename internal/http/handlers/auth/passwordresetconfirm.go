package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func PasswordResetConfirm(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email       string `json:"email"`
		ResetToken  string `json:"reset_token"`
		NewPassword string `json:"new_password"`
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

	// Check if the reset token matches for this user and email is verified
	var storedToken string
	err = db.QueryRow(
		"SELECT reset_token FROM users WHERE email = $1 AND email_verified = true",
		reqBody.Email,
	).Scan(&storedToken)
	if err != nil {
		http.Error(w, "User not found or email not verified", http.StatusBadRequest)
		log.Println("User not found or email not verified for:", reqBody.Email)
		return
	}

	if storedToken == "" || storedToken != reqBody.ResetToken {
		http.Error(w, "Invalid or expired reset token", http.StatusBadRequest)
		log.Println("Invalid reset token for:", reqBody.Email)
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		log.Println("Failed to hash password:", err)
		return
	}

	// Update the password and clear the reset token
	_, err = db.Exec(
		"UPDATE users SET password_hash = $1, reset_token = NULL, updated_at = NOW() WHERE email = $2",
		string(hashedPassword),
		reqBody.Email,
	)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		log.Println("Failed to update password for:", reqBody.Email, "Error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Password reset successfully",
		"data":    nil,
	}

	json.NewEncoder(w).Encode(response)

	log.Println("Password reset successfully for:", reqBody.Email)

}
