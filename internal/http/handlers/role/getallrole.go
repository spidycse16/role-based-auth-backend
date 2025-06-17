package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
	"github.com/sagorsarker04/Developer-Assignment/internal/utils"
)

func GetAllRole(w http.ResponseWriter, r *http.Request) {
	// Extract user info from the context
	userType := middleware.GetUserType(r)
	fmt.Println("User Type:", userType)

	// Allow only Admin and SystemAdmin
	if userType != "admin" && userType != "system_admin" {
		// http.Error(w, "No permission to access this resource", http.StatusForbidden)
		utils.ErrorResponse(w, http.StatusForbidden, "No permission to access this resource")
		return
	}

	// Connect to the database
	db := database.Connect()

	// Fetch all roles (now including description)
	rows, err := db.Query("SELECT id, name, description FROM roles ORDER BY created_at ASC")
	if err != nil {
		// http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch roles")
		return
	}
	defer rows.Close()

	// Collect roles
	var roles []map[string]interface{}
	for rows.Next() {
		var id, name, description string
		if err := rows.Scan(&id, &name, &description); err != nil {
			// http.Error(w, "Failed to read roles", http.StatusInternalServerError)
			utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to read roles")
			return
		}
		roles = append(roles, map[string]interface{}{
			"id":          id,
			"name":        name,
			"description": description,
		})
	}

	//Return the roles as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Roles retrieved successfully",
		"data":    roles,
	}

	json.NewEncoder(w).Encode(response)
	// utils.SuccessResponse(w, http.StatusOK, "Roles retrieved successfully", roles)
}
