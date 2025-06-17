package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
		"github.com/sagorsarker04/Developer-Assignment/internal/utils"
)

type DemoteRoleRequest struct {
	RoleName string `json:"role_name"`
}

func DemoteUserRole(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		// http.Error(w, "User ID is required", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Only system_admin and admin can demote
	userType := middleware.GetUserType(r)
	if userType != "system_admin" && userType != "admin" {
		// http.Error(w, "Only System Admin and Admin can demote users", http.StatusForbidden)
		utils.ErrorResponse(w, http.StatusForbidden, "Only System Admin and Admin can demote users")
		return
	}

	var req DemoteRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
		// http.Error(w, "Invalid request body", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Connect to the database
	db := database.Connect()

	var currentRole string
	err := db.QueryRow("SELECT user_type FROM users WHERE id = $1", userID).Scan(&currentRole)
	if err == sql.ErrNoRows {
		// http.Error(w, "User not found", http.StatusNotFound)
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	} else if err != nil {
		// http.Error(w, "Failed to fetch current role", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch current role")
		return
	}

	targetRole := "user"
	if req.RoleName != "" {
		targetRole = req.RoleName
	}

	if currentRole == "system_admin" {
		// http.Error(w, "System Admin cannot be demoted", http.StatusForbidden)
		utils.ErrorResponse(w, http.StatusForbidden, "System Admin cannot be demoted")
		return
	}

	if currentRole == "admin" {
		if targetRole != "moderator" && targetRole != "user" {
			// http.Error(w, "Admin can only be demoted to moderator or user", http.StatusBadRequest)
			utils.ErrorResponse(w, http.StatusBadRequest, "Admin can only be demoted to moderator or user")
			return
		}
	}

	if currentRole == "moderator" {
		if targetRole != "user" {
			// http.Error(w, "Moderator can only be demoted to user", http.StatusBadRequest)
			utils.ErrorResponse(w, http.StatusBadRequest, "Moderator can only be demoted to user")
			return
		}
	}

	if currentRole == "user" {
		// http.Error(w, "User cannot be demoted to any role", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "User cannot be demoted to any role")
		return
	}

	if targetRole == currentRole {
		// http.Error(w, "Target role must be different from current role", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "Target role must be different from current role")
		return
	}

	// Update users table
	_, err = db.Exec("UPDATE users SET user_type = $1 WHERE id = $2", targetRole, userID)
	if err != nil {
		// http.Error(w, "Cannot update user", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "Cannot update user")
		return
	}

	// Get role ID of the target role
	var newRoleID string
	err = db.QueryRow("SELECT id FROM roles WHERE name = $1", targetRole).Scan(&newRoleID)
	if err != nil {
		// http.Error(w, "Failed to fetch target role ID", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch target role ID")
		return
	}

	// Update user_roles table
	_, err = db.Exec("UPDATE user_roles SET role_id = $1 WHERE user_id = $2", newRoleID, userID)
	if err != nil {
		// http.Error(w, "Failed to update user role", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update user role")
		return
	}

	// Respond
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// response := map[string]interface{}{
	// 	"status":  strconv.Itoa(http.StatusOK),
	// 	"message": "User demoted successfully to " + targetRole,
	// 	"data":    nil,
	// }
	// json.NewEncoder(w).Encode(response)

	utils.SuccessResponse(w, http.StatusOK, "User demoted successfully to "+targetRole,nil)
}
