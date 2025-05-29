package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
	"github.com/sagorsarker04/Developer-Assignment/internal/utils"
)

type UpdateUserRequest struct {
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type UpdateUserResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	currentUserID := middleware.GetUserID(r)
	userType := middleware.GetUserType(r)

	// Only allow self-update or Admin/SystemAdmin
	if userID != currentUserID && userType != "admin" && userType != "system_admin" {
		// http.Error(w, "You are not authorized to update this user", http.StatusForbidden)
		utils.ErrorResponse(w, http.StatusForbidden, "You are not authorized to update this user")
		return
	}

	// Parse the request body
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// http.Error(w, "Invalid request payload", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request Payload")
		return
	}

	// Trim spaces from inputs
	req.Username = strings.TrimSpace(req.Username)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	// Connect to the database
	db := database.Connect()

	// Check if the new username is already taken (if provided)
	if req.Username != "" {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND id != $2)", req.Username, userID).Scan(&exists)
		if err != nil {
			// http.Error(w, "Failed to check username availability", http.StatusInternalServerError)
			utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to check username availability")
			return
		}
		if exists {
			// http.Error(w, "Username is already taken", http.StatusBadRequest)
			utils.ErrorResponse(w, http.StatusBadRequest, "Username is already taken")
			return
		}
	}

	// Build the update query directly
	query := `
	UPDATE users 
	SET 
		username = COALESCE(NULLIF($1, ''), username), 
		first_name = COALESCE(NULLIF($2, ''), first_name), 
		last_name = COALESCE(NULLIF($3, ''), last_name), 
		updated_at = NOW() 
	WHERE id = $4
	`

	// Execute the query
	_, err := db.Exec(query, req.Username, req.FirstName, req.LastName, userID)
	if err != nil {
		// http.Error(w, "Failed to update user", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "User updated successfully", nil)
}
