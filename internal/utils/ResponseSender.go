package utils

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse formats and sends a successful JSON response
func SuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  statusCode,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse formats and sends an error JSON response
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  statusCode,
		"message": message,
		"data":    nil,
	})
}
