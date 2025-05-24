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

// GetRoleDetails returns the details of a specific role
func GetRoleDetails(w http.ResponseWriter, r *http.Request) {
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
	db:=database.Connect()

	// Fetch role details
	var role struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE id = $1`
	err := db.QueryRow(query, roleID).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)

	// Handle not found and other errors
	if err == sql.ErrNoRows {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch role details", http.StatusInternalServerError)
		return
	}

	// Return the role details as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Role details retrieved successfully",
		"data":    role,
	}

	json.NewEncoder(w).Encode(response)

}
