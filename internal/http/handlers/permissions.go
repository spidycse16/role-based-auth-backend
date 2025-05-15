package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
	"github.com/gorilla/mux"
	"database/sql"
)

// ListAllPermissions lists all the permissions (Admin+)
func ListAllPermissions(w http.ResponseWriter, r *http.Request) {
	// Extract the user type from the context
	userType := middleware.GetUserType(r)

	// Allow only Admin and SystemAdmin
	if userType != "Admin" && userType != "SystemAdmin" {
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

	// Return the permissions as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permissions)
}

func GetPermissionDetails(w http.ResponseWriter, r *http.Request) {
	// Extract the user type from the context
	userType := middleware.GetUserType(r)

	// Allow only Admin and SystemAdmin
	if userType != "Admin" && userType != "SystemAdmin" {
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

	// Return the permission details as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permission)
}
