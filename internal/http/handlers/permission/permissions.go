package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
	"github.com/sagorsarker04/Developer-Assignment/internal/utils"
)

// ListAllPermissions lists all the permissions (Admin+)
func ListAllPermissions(w http.ResponseWriter, r *http.Request) {
	// Extract the user type from the context
	userType := middleware.GetUserType(r)

	// Allow only Admin and SystemAdmin
	if userType != "admin" && userType != "system_admin" {
		// http.Error(w, "No permission to access this resource", http.StatusForbidden)
		utils.ErrorResponse(w, http.StatusForbidden, "No permission to access this resource")
		return
	}

	// Connect to the database
	db := database.Connect()

	// Fetch all permissions
	rows, err := db.Query("SELECT id, name, description, resource, action, created_at, updated_at FROM permissions ORDER BY created_at ASC")
	if err != nil {
		// http.Error(w, "Failed to fetch permissions", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch permissions")
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
			// http.Error(w, "Failed to read permissions", http.StatusInternalServerError)
			utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to read permissions")
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

	// response := map[string]interface{}{
	// 	"status":  strconv.Itoa(http.StatusOK),
	// 	"message": "Permissions retrieved successfully",
	// 	"data":    permissions,
	// }

	// json.NewEncoder(w).Encode(response)
	utils.SuccessResponse(w, http.StatusOK, "Permissions retrieved successfully", permissions)

}

func GetPermissionDetails(w http.ResponseWriter, r *http.Request) {
	// Extract the user type from the context
	userType := middleware.GetUserType(r)

	// Allow only Admin and SystemAdmin
	if userType != "admin" && userType != "system_admin" {
		// http.Error(w, "No permission to access this resource", http.StatusForbidden)
		utils.ErrorResponse(w, http.StatusForbidden, "No permission to access this resource")
		return
	}

	// Get the permission ID from the URL
	vars := mux.Vars(r)
	permissionID := vars["permission_id"]

	// Connect to the database
	db := database.Connect()

	// Fetch the permission details
	var id, name, description, resource, action, createdAt, updatedAt string
	query := `
	SELECT id, name, description, resource, action, created_at, updated_at
	FROM permissions
	WHERE id = $1
	`
	err := db.QueryRow(query, permissionID).Scan(&id, &name, &description, &resource, &action, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		// http.Error(w, "Permission not found", http.StatusNotFound)
		utils.ErrorResponse(w, http.StatusNotFound, "Permission not found")
		return
	} else if err != nil {
		// http.Error(w, "Failed to fetch permission details", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch permission details")
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

	// json.NewEncoder(w).Encode(response)
	utils.SuccessResponse(w, http.StatusOK, "Permission details retrieved successfully", permission)

}
