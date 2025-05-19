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

// DeleteRole deletes a specific role
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	// Extract user info from the context
	userType := middleware.GetUserType(r)

	// Allow only Admin and SystemAdmin
	if userType != "admin" && userType != "system_admin" {
		http.Error(w, "No permission to access this resource", http.StatusForbidden)
		return
	}

	// Get the role_id from the URL path
	vars := mux.Vars(r)
	roleID := vars["role_id"]

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close(db)

	// Delete the role from the database
	query := "DELETE FROM roles WHERE id = $1 RETURNING id"
	var deletedRoleID string
	err = db.QueryRow(query, roleID).Scan(&deletedRoleID)

	if err == sql.ErrNoRows {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to delete role", http.StatusInternalServerError)
		return
	}

	// Return a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Role deleted successfully",
		"data":    nil,
	}

	json.NewEncoder(w).Encode(response)

}
