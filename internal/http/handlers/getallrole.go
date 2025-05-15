package handlers

import (
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

func GetAllRole(w http.ResponseWriter, r *http.Request) {
	// Extract user info from the context
	userType := middleware.GetUserType(r)
	fmt.Println(r)
	fmt.Println("User Type:", userType)
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

	// Fetch all roles
	rows, err := db.Query("SELECT id, name FROM roles ORDER BY created_at ASC")
	if err != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect roles
	var roles []map[string]interface{}
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			http.Error(w, "Failed to read roles", http.StatusInternalServerError)
			return
		}
		roles = append(roles, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	// Return the roles as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}
