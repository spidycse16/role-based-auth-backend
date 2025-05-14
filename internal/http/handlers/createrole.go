package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateRole(w http.ResponseWriter, r *http.Request) {
	// Extract user info from the context
	userType := middleware.GetUserType(r)

	// Allow only Admin and SystemAdmin
	if userType != "Admin" && userType != "SystemAdmin" {
		http.Error(w, "No permission to access this resource", http.StatusForbidden)
		return
	}

	// Decode the request body
	var req CreateRoleRequest
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

	// Insert the new role
	id := uuid.New()
	query := `
	INSERT INTO roles (id, name, description, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err = db.Exec(query, id, req.Name, req.Description, time.Now(), time.Now())
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "roles_name_key"` {
			http.Error(w, "Role name already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create role", http.StatusInternalServerError)
		}
		return
	}

	// Return the created role
	response := map[string]interface{}{
		"id":          id,
		"name":        req.Name,
		"description": req.Description,
		"created_at":  time.Now().Format(time.RFC3339),
		"updated_at":  time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
