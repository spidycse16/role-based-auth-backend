package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sagorsarker04/Developer-Assignment/internal/config"
	"github.com/sagorsarker04/Developer-Assignment/internal/utils"

	"github.com/sagorsarker04/Developer-Assignment/internal/database"
)

type ResetPasswordClaims struct {
	Email   string `json:"email"`
	Purpose string `json:"purpose"`
	jwt.RegisteredClaims
}

// PasswordResetRequest handles requests to initiate password reset by sending an email with a reset token.
func PasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email string `json:"email"`
	}

	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		// http.Error(w, "Invalid request body", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		log.Println("Invalid request body:", err)
		return
	}
	//frontend er jonno token
	cfg := config.GetConfig()
	secretKey := cfg.JWT.Secret

	claims := ResetPasswordClaims{
		Email:   reqBody.Email,
		Purpose: "password_reset",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		// http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		log.Println("JWT generation error:", err)
		return
	}

	fmt.Println(signedToken)
	//front er toekn upore

	db := database.Connect()

	//Generate a new reset token
	resetToken := uuid.NewString()
	email := reqBody.Email
	query := `Select id from users where email=$1`
	//Store the reset token in the database
	res, err := db.Exec(
		"UPDATE users SET reset_token = $1, updated_at = NOW() WHERE email = $2 AND email_verified = true",
		resetToken,
		reqBody.Email,
	)
	if err != nil {
		// http.Error(w, "Failed to store reset token", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to store reset token")
		log.Println("Failed to store reset token for:", reqBody.Email, "Error:", err)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		// http.Error(w, "No verified user found with this email", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "No verified user found with this email")
		log.Println("No verified user found with email:", reqBody.Email)
		return
	}
	var userID int
	err = db.QueryRow(query, email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No user found with this email")
		} else {
			fmt.Println("Invalid query")
		}
	}
	fmt.Println("User id is", userID)

	if err := sendResetToken(reqBody.Email, resetToken, signedToken); err != nil {
		// http.Error(w, "Failed to send password reset email", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to send password reset email")
		log.Println("Failed to send password reset email to:", reqBody.Email, "Error:", err)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Password reset email sent successfully", nil)

}

// sendResetToken sends a password reset email to the specified email address.
func sendResetToken(toEmail, resetToken, signedToken string) error {
	cfg := config.GetConfig()

	from := cfg.Email.From
	password := cfg.Email.Password
	smtpHost := cfg.Email.Host
	smtpPort := cfg.Email.Port

	subject := "Password Reset Request"

	resetLink := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", signedToken)

	// Use HTML content for email body
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<h2>Password Reset Request</h2>
			<p>Hello,</p>
			<p>You requested to reset your password. Use the following verification code:</p>
			<h3 style="background-color: #f3f3f3; padding: 10px; display: inline-block;">%s</h3>
			<p>Or click the button below to reset your password directly:</p>
			<a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #007BFF; color: white; text-decoration: none; border-radius: 4px;">Reset Password</a>
			<p>If you did not request a password reset, please ignore this email.</p>
			<p>Thank you.</p>
		</body>
		</html>
	`, resetToken, resetLink)

	// Add HTML content-type in headers
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		from, toEmail, subject, body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", smtpHost, smtpPort),
		auth,
		from,
		[]string{toEmail},
		[]byte(message),
	)
	if err != nil {
		log.Println("Failed to send password reset email to:", toEmail, "Error:", err)
		return err
	}

	log.Println("Password reset email sent successfully to:", toEmail)
	return nil
}
