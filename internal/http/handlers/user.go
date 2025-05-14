package handlers

import (
	// "encoding/json"
	"net/http"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) string {
	userID := middleware.GetUserID(r)
	// username := middleware.GetUsername(r)
	// userType := middleware.GetUserType(r)

	// Return the extracted user info
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"user_id":   userID,
	// 	"username":  username,
	// 	"user_type": userType,
	// })
	return userID
}
