package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// LogoutUser handles the user logout process by clearing the auth cookie.
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Clear the auth_token cookie
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "", // Clear the value
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // Set expiry to the past
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}

	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Logout successful",
		"data":    nil,
	}

	json.NewEncoder(w).Encode(response)

}
