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

// ListAllPermissions lists all the permissions (Admin+)
func ListAllPermissions(w http.ResponseWriter, r *http.Request) {
	// Extract the user type from the context
	userType := middleware.GetUserType(r)

	// Allow only Admin and SystemAdmin
	if userType != "admin" && userType != "system_admin" {
		http.Error(w, "No permission to access this resource", http.StatusForbidden)
		return
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close(db)

	// Fetch all permissions
	rows, err := db.Query("SELECT id, name, description, resource, action, created_at, updated_at FROM permissions ORDER BY created_at ASC")
	if err != nil {
		http.Error(w, "Failed to fetch permissions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect permissions
	var permissions []map[string]interface{}
	for rows.Next() {
		var id, name, description, resource, action string
		var createdAt, updatedAt string

		err := rows.Scan(&id, &name, &description, &resource, &action, &createdAt, &updatedAt)
		if err != nil {
			http.Error(w, "Failed to read permissions", http.StatusInternalServerError)
			return
		}

		permissions = append(permissions, map[string]interface{}{
			"id":          id,
			"name":        name,
			"description": description,
			"resource":    resource,
			"action":      action,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Permissions retrieved successfully",
		"data":    permissions,
	}

	json.NewEncoder(w).Encode(response)

}

func GetPermissionDetails(w http.ResponseWriter, r *http.Request) {
	// Extract the user type from the context
	userType := middleware.GetUserType(r)

	// Allow only Admin and SystemAdmin
	if userType != "admin" && userType != "system_admin" {
		http.Error(w, "No permission to access this resource", http.StatusForbidden)
		return
	}

	// Get the permission ID from the URL
	vars := mux.Vars(r)
	permissionID := vars["permission_id"]

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close(db)

	// Fetch the permission details
	var id, name, description, resource, action, createdAt, updatedAt string
	query := `
	SELECT id, name, description, resource, action, created_at, updated_at
	FROM permissions
	WHERE id = $1
	`
	err = db.QueryRow(query, permissionID).Scan(&id, &name, &description, &resource, &action, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Permission not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch permission details", http.StatusInternalServerError)
		return
	}

	// Build the response
	permission := map[string]interface{}{
		"id":          id,
		"name":        name,
		"description": description,
		"resource":    resource,
		"action":      action,
		"created_at":  createdAt,
		"updated_at":  updatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Permission details retrieved successfully",
		"data":    permission,
	}

	json.NewEncoder(w).Encode(response)

}
