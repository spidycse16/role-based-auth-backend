package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
	"github.com/sagorsarker04/Developer-Assignment/internal/models"
)



// CreatePermission handles creating a new permission
func CreatePermission(w http.ResponseWriter, r *http.Request) {
	// Extract user info from the context
	userType := middleware.GetUserType(r)

	// Allow only Admin and SystemAdmin
	if userType != "Admin" && userType != "SystemAdmin" {
		http.Error(w, "No permission to access this resource", http.StatusForbidden)
		return
	}

	// Parse the request body
	var req models.PermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" || req.Resource == "" || req.Action == "" {
		http.Error(w, "Name, resource, and action are required", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close(db)

	// Insert the permission into the database
	query := `
	INSERT INTO permissions (id, name, description, resource, action, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	permissionID := uuid.NewString()
	createdAt := time.Now()

	_, err = db.Exec(query, permissionID, req.Name, req.Description, req.Resource, req.Action, createdAt, createdAt)
	if err != nil {
		http.Error(w, "Failed to create permission", http.StatusInternalServerError)
		return
	}

	// Return the created permission ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":        "Permission created successfully",
		"permission_id":  permissionID,
	})
}
