package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

type DemoteRoleRequest struct {
	RoleName string `json:"role_name"`
}

func DemoteUserRole(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Protect System Admin from being demoted
	userType := middleware.GetUserType(r)
	if userType != "system_admin" && userType != "admin" {
		http.Error(w, "Only System Admin and Admin can demote users", http.StatusForbidden)
		return
	}

	// Parse the optional role name from the request body
	var req DemoteRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close(db)

	// Get the current role of the target user
	var currentRole string
	err = db.QueryRow("SELECT user_type FROM users WHERE id = $1", userID).Scan(&currentRole)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch current role", http.StatusInternalServerError)
		return
	}

	// Determine the target role index
	targetRole := "user" // Default to demoting to User if no role is provided
	if req.RoleName != "" {
		targetRole = req.RoleName
	}

	// Validate demotion logic with if-else
	if currentRole == "system_admin" {
		// Already handled earlier, but just for clarity
		http.Error(w, "System Admin cannot be demoted", http.StatusForbidden)
		return
	}

	if currentRole == "admin" {
		// Admin can only be demoted to moderator or user
		if targetRole != "moderator" && targetRole != "user" {
			http.Error(w, "Admin can only be demoted to moderator or user", http.StatusBadRequest)
			return
		}
	}

	if currentRole == "moderator" {
		// Moderator cannot demote or be demoted to higher roles, but this is demote API so:
		// Demoting moderator only allowed to user
		if targetRole != "user" {
			http.Error(w, "Moderator can only be demoted to user", http.StatusBadRequest)
			return
		}
	}

	if currentRole == "user" {
		// User cannot be demoted to anything else
		http.Error(w, "User cannot be demoted to any role", http.StatusBadRequest)
		return
	}

	// If targetRole is the same as currentRole, reject
	if targetRole == currentRole {
		http.Error(w, "Target role must be different from current role", http.StatusBadRequest)
		return
	}

	// Perform the update in DB
	_, err = db.Exec("UPDATE users SET user_type = $1 WHERE id = $2", targetRole, userID)
	if err != nil {
		http.Error(w, "Cannot update user", http.StatusBadRequest)
		return
	}

	// Return a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "User demoted successfully to " + targetRole,
		"data":    nil,
	}

	json.NewEncoder(w).Encode(response)

}
