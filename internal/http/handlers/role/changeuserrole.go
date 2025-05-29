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

type ChangeRoleRequest struct {
	RoleName string `json:"role_name"`
}

func ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	currentUserType := middleware.GetUserType(r)

	if userID == "" {
		// http.Error(w, "User ID is required", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}

	if currentUserType != "admin" && currentUserType != "system_admin" {
		// http.Error(w, "Only Admin or System Admin can change the role", http.StatusForbidden)
		utils.ErrorResponse(w, http.StatusForbidden, "Only Admin or System Admin can change the role")
		return
	}

	var req ChangeRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RoleName == "" {
		// http.Error(w, "Invalid request body", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Connect to the database
	db := database.Connect()

	// Fetch the current user type
	var currentRole string
	err := db.QueryRow("SELECT user_type FROM users WHERE id = $1", userID).Scan(&currentRole)
	if err == sql.ErrNoRows {
		// http.Error(w, "User not found", http.StatusNotFound)
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	} else if err != nil {
		// http.Error(w, "Failed to fetch current user role", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch current user role")
		return
	}

	// Prevent promotion of regular users to Admin, SystemAdmin, or Moderator
	if currentRole == "user" && (req.RoleName == "admin" || req.RoleName == "system_admin" || req.RoleName == "moderator") {
		// http.Error(w, "Regular users cannot be promoted to Admin, SystemAdmin, or Moderator", http.StatusForbidden)
		utils.ErrorResponse(w, http.StatusForbidden, "Regular users cannot be promoted to Admin, SystemAdmin, or Moderator")
		return
	}

	// Find the role ID for the given role name
	var roleID string
	err = db.QueryRow("SELECT id FROM roles WHERE name = $1", req.RoleName).Scan(&roleID)
	if err == sql.ErrNoRows {
		// http.Error(w, "Role not found", http.StatusNotFound)
		utils.ErrorResponse(w, http.StatusNotFound, "Role not found")
		return
	} else if err != nil {
		// http.Error(w, "Failed to fetch role", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch role")
		return
	}

	// Update the user's main role in the users table
	_, err = db.Exec("UPDATE users SET user_type = $1, updated_at = NOW() WHERE id = $2", req.RoleName, userID)
	if err != nil {
		// http.Error(w, "Failed to update user type", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update user type")
		return
	}

	// Update the user's role in the user_roles table
	_, err = db.Exec("UPDATE user_roles SET role_id = $1 WHERE user_id = $2", roleID, userID)
	if err != nil {
		// http.Error(w, "Failed to update user role", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update user role")
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	// response := map[string]interface{}{
	// 	"status":  strconv.Itoa(http.StatusOK),
	// 	"message": "User role updated successfully",
	// 	"data":    nil,
	// }

	// json.NewEncoder(w).Encode(response)

	utils.SuccessResponse(w, http.StatusOK, "User role updated successfully", nil)

}
