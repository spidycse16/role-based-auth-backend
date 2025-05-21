package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sagorsarker04/Developer-Assignment/internal/config"
)

// GetCurrentUserProfile returns the current authenticated user's profile
func GetCurrentUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get the cookie
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "No valid authentication token", http.StatusUnauthorized)
		return
	}
	
	// Load the config
	cfg:=config.GetConfig()
	// if err != nil {
	// 	http.Error(w, "Failed to load config", http.StatusInternalServerError)
	// 	return
	// }

	// Parse the JWT token
	tokenString := cookie.Value
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Extract the token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID := claims["user_id"]
		username := claims["username"]
		userType := claims["user_type"]
		exp := claims["exp"]

		expirationTime := time.Unix(int64(exp.(float64)), 0)

		// Build the response
		response := map[string]interface{}{
			"user_id":   userID,
			"username":  username,
			"user_type": userType,
			"expires_at": expirationTime.Format(time.RFC3339),
			"who_am_i":  userType,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response["who_am_i"])
	} else {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
	}
}
