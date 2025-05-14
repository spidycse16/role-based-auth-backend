package middleware

import (
	"net/http"
)

// UserContextKeys defines the keys for user info stored in the context
type UserContextKeys string

const (
	UserIDKey   UserContextKeys = "user_id"
	UsernameKey UserContextKeys = "username"
	UserTypeKey UserContextKeys = "user_type"
)

// GetUserID extracts the user ID from the request context
func GetUserID(r *http.Request) string {
	if userID, ok := r.Context().Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetUsername extracts the username from the request context
func GetUsername(r *http.Request) string {
	if username, ok := r.Context().Value(UsernameKey).(string); ok {
		return username
	}
	return ""
}

// GetUserType extracts the user type from the request context
func GetUserType(r *http.Request) string {
	if userType, ok := r.Context().Value(UserTypeKey).(string); ok {
		return userType
	}
	return ""
}
