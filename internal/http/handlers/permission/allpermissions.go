package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

// GetCurrentUserPermissions returns the permissions for the authenticated user
func GetCurrentUserPermissions(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the context
	userID := middleware.GetUserID(r)

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Fprintln(w, userID)
	db:=database.Connect()

	// Fetch the user's permissions
	query := `
	SELECT p.id, p.name, p.description, p.resource, p.action
	FROM user_roles ur
	INNER JOIN role_permissions rp ON ur.role_id = rp.role_id
	INNER JOIN permissions p ON rp.permission_id = p.id
	WHERE ur.user_id = $1
	ORDER BY p.name ASC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		http.Error(w, "Failed to fetch permissions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect permissions
	var permissions []map[string]interface{}
	for rows.Next() {
		var id, name, description, resource, action string
		if err := rows.Scan(&id, &name, &description, &resource, &action); err != nil {
			http.Error(w, "Failed to read permissions", http.StatusInternalServerError)
			return
		}
		permissions = append(permissions, map[string]interface{}{
			"id":          id,
			"name":        name,
			"description": description,
			"resource":    resource,
			"action":      action,
		})
	}

	// Return the permissions as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Permissions fetched successfully",
		"data":    permissions,
	}

	json.NewEncoder(w).Encode(response)

}
