package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
	"github.com/sagorsarker04/Developer-Assignment/internal/utils"
)

// DeleteRole deletes a specific role
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	// Extract user info from the context
	userType := middleware.GetUserType(r)

	// Allow only Admin and SystemAdmin
	if userType != "admin" && userType != "system_admin" {
		// http.Error(w, "No permission to access this resource", http.StatusForbidden)
		utils.ErrorResponse(w, http.StatusForbidden, "No permission to access this resource")
		return
	}

	// Get the role_id from the URL path
	vars := mux.Vars(r)
	roleID := vars["role_id"]

	// Connect to the database
	db := database.Connect()

	// Delete the role from the database
	query := "DELETE FROM roles WHERE id = $1 RETURNING id"
	var deletedRoleID string
	err := db.QueryRow(query, roleID).Scan(&deletedRoleID)

	if err == sql.ErrNoRows {
		// http.Error(w, "Role not found", http.StatusNotFound)
		utils.ErrorResponse(w, http.StatusNotFound, "Role not found")
		return
	} else if err != nil {
		// http.Error(w, "Failed to delete role", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete role")
		return
	}

	// Return a success message

	// json.NewEncoder(w).Encode(response)
	utils.SuccessResponse(w, http.StatusOK, "Role deleted successfully", nil)

}
