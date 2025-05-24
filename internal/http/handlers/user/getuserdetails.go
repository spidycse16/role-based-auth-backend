package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	// Extract user type and user ID from the context
	userType := middleware.GetUserType(r)
	userID := middleware.GetUserID(r)
	fmt.Println("Extracted User Type:", userType)
	fmt.Println("Extracted User ID:", userID)

	//Allow only authenticated users
	if userType == "" || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the requested user ID from the URL
	requestedUserID := mux.Vars(r)["user_id"]

	//Allow only the user themselves or Admin/SystemAdmin
	if requestedUserID != userID && userType != "admin" && userType != "system_admin" {
		http.Error(w, "No permission to access this resource", http.StatusForbidden)
		return
	}

	// Connect to the database
	db:=database.Connect()
	
	// Fetch the user details
	var user map[string]interface{}
	query := `
		SELECT id, username, email, first_name, last_name, user_type, email_verified, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	row := db.QueryRow(query, requestedUserID)

	var id, username, email, firstName, lastName, userTypeDB string
	var emailVerified bool
	var createdAt, updatedAt string

	if err := row.Scan(&id, &username, &email, &firstName, &lastName, &userTypeDB, &emailVerified, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch user details", http.StatusInternalServerError)
		}
		return
	}

	user = map[string]interface{}{
		"id":             id,
		"username":       username,
		"email":          email,
		"first_name":     firstName,
		"last_name":      lastName,
		"user_type":      userTypeDB,
		"email_verified": emailVerified,
		"created_at":     createdAt,
		"updated_at":     updatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "User details retrieved successfully",
		"data":    user,
	}

	json.NewEncoder(w).Encode(response)

}
