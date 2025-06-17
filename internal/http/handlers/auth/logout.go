package handlers

import (
	"net/http"

	"time"
	"github.com/sagorsarker04/Developer-Assignment/internal/utils"
)

// LogoutUser handles the user logout process by clearing the auth cookie.
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Clear the auth_token cookie
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "", 
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), 
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(w, cookie)

	// json.NewEncoder(w).Encode(response)
	utils.SuccessResponse(w, http.StatusOK, "Logout successful", nil)

}
