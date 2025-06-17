package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sagorsarker04/Developer-Assignment/internal/config"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

//var jwtSecret = []byte("your_secret_key") // Replace with your actual secret

type Claims struct {
	Email   string `json:"email"`
	Purpose string `json:"purpose"`
	jwt.RegisteredClaims
}

func PasswordReset(w http.ResponseWriter, r *http.Request) {

	cfg:=config.GetConfig()
	jwtSecret:=cfg.JWT.Secret

	type RequestBody struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		// http.Error(w, "Invalid request body", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		log.Println("Invalid request body:", err)
		return
	}
	if len(reqBody.NewPassword) < 8 {
	utils.ErrorResponse(w, http.StatusBadRequest, "Password must be at least 8 characters")
	return
}

	// Parse and validate JWT token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(reqBody.Token, claims, func(token *jwt.Token) (interface{}, error) {
    return []byte(jwtSecret), nil 
})
	if err != nil || !token.Valid {
		// http.Error(w, "Invalid or expired reset token", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid or expired reset token")
		log.Println("Invalid or expired token:", err)
		return
	}

	if claims.Purpose != "password_reset" {
		// http.Error(w, "Invalid token purpose", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid token purpose")
		log.Println("Invalid token purpose for email:", claims.Email)
		return
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		// http.Error(w, "Token has expired", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "Token has expired")
		log.Println("Token expired for email:", claims.Email)
		return
	}

	// Check if the user exists and email is verified
	db := database.Connect()
	var emailVerified bool
	err = db.QueryRow(
		"SELECT email_verified FROM users WHERE email = $1",
		claims.Email,
	).Scan(&emailVerified)
	if err != nil || !emailVerified {
		// http.Error(w, "User not found or email not verified", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "User not found or email not verified")
		log.Println("User not found or email not verified for:", claims.Email)
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		// http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to hash password")
		log.Println("Failed to hash password:", err)
		return
	}

	// Update the password
	_, err = db.Exec(
		"UPDATE users SET password_hash = $1, updated_at = NOW() WHERE email = $2",
		string(hashedPassword),
		claims.Email,
	)
	if err != nil {
		// http.Error(w, "Failed to update password", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update password")
		log.Println("Failed to update password for:", claims.Email, "Error:", err)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Password reset successfully", nil)
	log.Println("Password reset successfully for:", claims.Email)
}
