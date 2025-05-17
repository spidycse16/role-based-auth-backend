package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
	"github.com/sagorsarker04/Developer-Assignment/internal/models"
)



// UpdateRole updates the details of a specific role
func UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Extract user info from the context
	userType := middleware.GetUserType(r)

	// Allow only Admin and SystemAdmin
	if userType != "Admin" && userType != "SystemAdmin" {
		http.Error(w, "No permission to access this resource", http.StatusForbidden)
		return
	}

	// Get the role_id from the URL path
	vars := mux.Vars(r)
	roleID := vars["role_id"]

	// Parse the request body
	var req models.RoleUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		http.Error(w, "Role name is required", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close(db)

	// Update the role in the database
	query := `
	UPDATE roles
	SET name = $1, description = $2, updated_at = NOW()
	WHERE id = $3
	RETURNING id, name, description, created_at, updated_at
	`
	var updatedRole struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}

	err = db.QueryRow(query, req.Name, req.Description, roleID).Scan(
		&updatedRole.ID,
		&updatedRole.Name,
		&updatedRole.Description,
		&updatedRole.CreatedAt,
		&updatedRole.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to update role", http.StatusInternalServerError)
		return
	}

	// Return the updated role as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedRole)
}
