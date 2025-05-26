package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"

	"github.com/google/uuid"
	"github.com/sagorsarker04/Developer-Assignment/internal/config"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
)

// PasswordResetRequest handles requests to initiate password reset by sending an email with a reset token.
func PasswordResetRequest(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to store reset token", http.StatusInternalServerError)
		log.Println("Failed to store reset token for:", reqBody.Email, "Error:", err)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "No verified user found with this email", http.StatusBadRequest)
		log.Println("No verified user found with email:", reqBody.Email)
		return
	}
	var userID int
	err = db.QueryRow(query,email).Scan(&userID)
	if err!=nil{
		if err==sql.ErrNoRows{
			fmt.Println("No user found with this email")
		}else{
			fmt.Println("Invalid query")
		}
	}
	fmt.Println("User id is",userID)

	// email := reqBody.Email

	// var userID uuid.UUID
	// query := `SELECT id FROM users WHERE email = $1 AND email_verified = true`

	// err := db.QueryRow(query, email).Scan(&userID)
	// fmt.Println("user id is",userID)
	// if err == sql.ErrNoRows {
	// 	fmt.Println("No verified user found with this email")
	// 	return
	// } else if err != nil {
	// 	fmt.Println("Query error:", err)
	// 	return
	// }

	// //fmt.Println("User id is", userID)



	// resetToken := "dfuhfuh"

	// Send the password reset email
	if err := sendResetToken(reqBody.Email, resetToken); err != nil {
		http.Error(w, "Failed to send password reset email", http.StatusInternalServerError)
		log.Println("Failed to send password reset email to:", reqBody.Email, "Error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Password reset email sent successfully",
		"data":    nil,
	}

	json.NewEncoder(w).Encode(response)

	log.Println("Password reset email sent successfully to:", reqBody.Email)

}

// sendResetToken sends a password reset email to the specified email address.
func sendResetToken(toEmail, resetToken string) error {
	cfg := config.GetConfig()
	// if err != nil {
	// 	log.Println("Failed to load config:", err)
	// 	return err
	// }

	from := cfg.Email.From
	password := cfg.Email.Password
	smtpHost := cfg.Email.Host
	smtpPort := cfg.Email.Port

	subject := "Password Reset Verification Code"
	body := fmt.Sprintf(
		"Hello,\n\nYou requested to reset your password. Use the following verification code to reset it:\n\n%s\n\nIf you did not request a password reset, please ignore this email.\n\nThank you.",
		resetToken,
	)
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, toEmail, subject, body)

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

	log.Println("Password reset verification code sent successfully to:", toEmail)
	return nil
}
